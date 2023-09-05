package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/framework/gin-gonic/hfwctx"
	"github.com/charlesbases/library/framework/gin-gonic/webserver"
)

// upgrader websocker upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: time.Second * 3,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

type sessionID string

// Session .
type Session struct {
	id            sessionID
	subscriptions map[subject]bool

	// request 请求
	request chan *WebSocketRequest
	// response 相应
	response chan *WebSocketResponse
	// broadcast 广播推送
	broadcast chan *WebSocketBroadcast

	ctx  *hfwctx.Context
	conn *websocket.Conn

	opts *Options

	lock sync.RWMutex
	once sync.Once

	ready   bool
	closed  bool
	closing chan struct{}
}

// ping .
func (session *Session) ping() {
	if !session.closed {
		session.ready = true

		session.writeEvent(MethodPing, library.NowString())
	}
}

// serve .
func (session *Session) serve() {
	go session.listening()

	for {
		select {
		case <-time.After(session.opts.Heartbeat):
			if session.ready {
				session.ready = false
				session.ctx.Debugf("[WebSocketID: %d] no heartbeat", session.id)
			}
		case request, ok := <-session.request:
			if ok {
				switch request.Method {
				case MethodPing:
					session.ping()
				case MethodSubscribe:
					session.subscribe(request.Params)
				case MethodUnsubscribe:
					session.unsubscribe(request.Params)
				case MethodDisconnect:
					session.close()
				default:
					session.WriteError(webserver.StatusParamInvalid.WebError(fmt.Sprintf("invalid method: %s", request.Method)))
				}
			}
		case response, ok := <-session.response:
			if ok {
				session.write(response)
			}
		case event, ok := <-session.broadcast:
			if ok {
				session.event(event)
			}
		case <-session.closing:
			session.disconnect()
			return
		}
	}
}

// listening .
func (session *Session) listening() {
	for {
		select {
		case <-session.closing:
			return
		default:
			request := new(WebSocketRequest)
			if err := session.conn.ReadJSON(request); err != nil {
				switch session.isCloseError(err) {
				case websocket.CloseNormalClosure, websocket.CloseNoStatusReceived:
				default:
					session.ctx.Errorf("[WebSocketID: %d] read message error: %s", session.id, err.Error())
					session.WriteError(webserver.StatusBadRequest.WebError(err.Error()))
				}

				session.close()
				return
			}

			if !store.verifySession(request.ID) {
				session.WriteError(webserver.StatusParamInvalid.WebError(fmt.Sprintf("invalid session id, %d not connected", request.ID)))
				session.close()
				return
			}

			if session.closed {
				session.WriteError(webserver.StatusBadRequest.WebError("connect closed"))
				return
			}

			if request.Method == MethodPing {
				goto r
			}

			if !session.ready {
				session.WriteError(webserver.StatusBadRequest.WebError("connect not ready"))
				continue
			}

			if request.Params == nil {
				session.WriteError(webserver.StatusParamInvalid.WebError("params cannot be empty."))
				continue
			}

		r:
			session.ctx.Debugf("[WebSocketID: %d] [r] Method: %s", session.id, request.Method.String())
			session.request <- request
		}
	}
}

// isCloseError .
func (session *Session) isCloseError(err error) int {
	if e, ok := err.(*websocket.CloseError); ok {
		return e.Code
	} else {
		return -1
	}
}

// write .
func (session *Session) write(v *WebSocketResponse) error {
	if session.ready {
		v.ID = session.id
		session.ctx.Debugf("[WebSocketID: %d] [w] Method: %s", session.id, v.Method.String())
		return session.conn.WriteJSON(v)
	}

	session.ctx.Debugf("[WebSocketID: %d] write failed. connect not ready.", session.id)
	return nil
}

// writeEvent WebSocket 返回 Method 相关信息
func (session *Session) writeEvent(method Method, data interface{}) {
	session.response <- &WebSocketResponse{
		ID:     session.id,
		Method: method,
		Data:   data,
	}
}

// WriteResponse 手动推送 response 至 WebSocket
func (session *Session) WriteResponse(data interface{}) {
	session.writeEvent(MethodResponse, data)
}

// WriteError WebSocket 返回错误
func (session *Session) WriteError(err *webserver.WebError) {
	session.response <- &WebSocketResponse{
		ID:      session.id,
		Code:    err.Code,
		Message: err.Message,
	}
}

// Context .
func (session *Session) Context() *hfwctx.Context {
	return session.ctx
}

type subject string

// verify .
func (sub subject) verify(tar subject) bool {
	if sub == "*" {
		return true
	}
	if sub == tar {
		return true
	}
	if strings.HasSuffix(string(sub), "*") {
		prefix := strings.TrimSuffix(string(sub), "*")
		return strings.HasPrefix(string(tar), prefix)
	}
	return false
}

// event .
func (session *Session) event(event *WebSocketBroadcast) {
	if session.ready {
		session.lock.RLock()
		for subject := range session.subscriptions {
			if subject.verify(event.Subject) {
				session.writeEvent(MethodBroadcast, event)
				break
			}
		}
		session.lock.RUnlock()
	}
}

// topics .
func (session *Session) topics() []subject {
	session.lock.RLock()

	var topics = make([]subject, 0, len(session.subscriptions))
	for subject := range session.subscriptions {
		topics = append(topics, subject)
	}
	session.lock.RUnlock()

	return topics
}

// subscribe .
func (session *Session) subscribe(params *json.RawMessage) {
	var subjects = make([]subject, 0)

	if err := json.Unmarshal(*params, &subjects); err != nil {
		session.ctx.Errorf("[WebSocketID: %d] subscribe failed. %s", session.id, err.Error())
		session.WriteError(webserver.StatusParamInvalid.WebError(err.Error()))
		session.close()
		return
	}

	var subscribers = make([]*subscriber, 0, len(subjects))

	session.lock.Lock()
	for _, subject := range subjects {
		session.subscriptions[subject] = true

		subscribers = append(subscribers, &subscriber{
			sessionID: session.id,
			subject:   subject,
			onEvent:   session.broadcast,
		})
	}
	session.lock.Unlock()

	go es.subscribe(subscribers...)

	session.writeEvent(MethodSubscribe, session.topics())
}

// unsubscribe .
func (session *Session) unsubscribe(params *json.RawMessage) {
	var subjects = make([]subject, 0)

	if err := json.Unmarshal(*params, &subjects); err != nil {
		session.ctx.Errorf("[WebSocketID: %d] unsubscribe failed. %s", session.id, err.Error())
		session.WriteError(webserver.StatusParamInvalid.WebError(err.Error()))
		session.close()
		return
	}

	var subscribers = make([]*subscriber, 0, len(subjects))

	session.lock.Lock()
	for _, subject := range subjects {
		delete(session.subscriptions, subject)

		subscribers = append(subscribers, &subscriber{
			sessionID: session.id,
			subject:   subject,
			onEvent:   session.broadcast,
		})
	}
	session.lock.Unlock()

	go es.unsubscribe(subscribers...)

	session.writeEvent(MethodUnsubscribe, session.topics())
}

// close exec session.disconnect
func (session *Session) close() {
	session.closing <- struct{}{}
}

// disconnect .
func (session *Session) disconnect() {
	session.once.Do(func() {
		session.ctx.Debugf("[WebSocketID: %d] disconnected", session.id)

		session.ready = false
		session.closed = true

		store.dropSession(session.id)

		close(session.closing)

		close(session.request)
		close(session.response)
		close(session.broadcast)

		if session.conn != nil {
			session.conn.Close()
			session.conn = nil
		}

		session = nil
	})
}

var store = &pool{store: make(map[sessionID]struct{}, 0)}

// pool .
type pool struct {
	lk    sync.RWMutex
	store map[sessionID]struct{}
}

// verifySession .
func (pool *pool) verifySession(id sessionID) bool {
	pool.lk.RLock()
	_, found := pool.store[id]
	pool.lk.RUnlock()

	return found
}

// newSession .
func (pool *pool) createSession() sessionID {
	id := sessionID(uuid.NewString())

	pool.lk.Lock()
	pool.store[id] = struct{}{}
	pool.lk.Unlock()

	return id
}

// dropSession .
func (pool *pool) dropSession(id sessionID) {
	pool.lk.Lock()
	delete(pool.store, id)
	pool.lk.Unlock()
}
