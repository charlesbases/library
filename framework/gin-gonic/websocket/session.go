package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/charlesbases/logger"
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

// Session .
type Session struct {
	id            hfwctx.ID
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
func (sess *Session) ping() {
	if !sess.closed {
		sess.ready = true

		sess.writeEvent(MethodPing, library.NowString())
	}
}

// serve .
func (sess *Session) serve() {
	go sess.listening()

	for {
		select {
		case <-time.After(sess.opts.Heartbeat):
			if sess.ready {
				sess.ready = false
				logger.Context(sess.ctx).Debugf("[WebSocketID: %s] no heartbeat", sess.id)
			}
		case request, ok := <-sess.request:
			if ok {
				switch request.Method {
				case MethodPing:
					sess.ping()
				case MethodSubscribe:
					sess.subscribe(request.Params)
				case MethodUnsubscribe:
					sess.unsubscribe(request.Params)
				case MethodDisconnect:
					sess.close()
				default:
					sess.WriteError(webserver.StatusParamInvalid.WebError(fmt.Sprintf("invalid method: %s", request.Method)))
				}
			}
		case response, ok := <-sess.response:
			if ok {
				sess.write(response)
			}
		case event, ok := <-sess.broadcast:
			if ok {
				sess.event(event)
			}
		case <-sess.closing:
			sess.disconnect()
			return
		}
	}
}

// listening .
func (sess *Session) listening() {
	for {
		select {
		case <-sess.closing:
			return
		default:
			request := new(WebSocketRequest)
			if err := sess.conn.ReadJSON(request); err != nil {
				switch sess.isCloseError(err) {
				case websocket.CloseNormalClosure, websocket.CloseNoStatusReceived:
				default:
					logger.Context(sess.ctx).Errorf("[WebSocketID: %s] read message error: %s", sess.id, err)
					sess.WriteError(webserver.StatusBadRequest.WebError(err))
				}

				sess.close()
				return
			}

			if !store.verifySession(request.ID) {
				sess.WriteError(
					webserver.StatusParamInvalid.WebError(
						fmt.Sprintf(
							"invalid sess id, %s not connected", request.ID,
						),
					),
				)
				sess.close()
				return
			}

			if sess.closed {
				sess.WriteError(webserver.StatusBadRequest.WebError("connect closed"))
				return
			}

			if request.Method == MethodPing {
				goto r
			}

			if !sess.ready {
				sess.WriteError(webserver.StatusBadRequest.WebError("connect not ready"))
				continue
			}

			if request.Params == nil {
				sess.WriteError(webserver.StatusParamInvalid.WebError("params cannot be empty."))
				continue
			}

		r:
			logger.Context(sess.ctx).Debugf("[WebSocketID: %s] [r] Method: %s", sess.id, request.Method.String())
			sess.request <- request
		}
	}
}

// isCloseError .
func (sess *Session) isCloseError(err error) int {
	if e, ok := err.(*websocket.CloseError); ok {
		return e.Code
	} else {
		return -1
	}
}

// write .
func (sess *Session) write(v *WebSocketResponse) error {
	if sess.ready {
		v.ID = sess.id
		logger.Context(sess.ctx).Debugf("[WebSocketID: %s] [w] Method: %s", sess.id, v.Method.String())
		return sess.conn.WriteJSON(v)
	}

	logger.Context(sess.ctx).Debugf("[WebSocketID: %s] write failed. connect not ready.", sess.id)
	return nil
}

// writeEvent WebSocket 返回 Method 相关信息
func (sess *Session) writeEvent(method Method, data interface{}) {
	sess.response <- &WebSocketResponse{
		ID:     sess.id,
		Method: method,
		Data:   data,
	}
}

// WriteResponse 手动推送 response 至 WebSocket
func (sess *Session) WriteResponse(data interface{}) {
	sess.writeEvent(MethodResponse, data)
}

// WriteError WebSocket 返回错误
func (sess *Session) WriteError(err *webserver.WebError) {
	sess.response <- &WebSocketResponse{
		ID:      sess.id,
		Code:    err.Code,
		Message: err.Message,
	}
}

// Context .
func (sess *Session) Context() *hfwctx.Context {
	return sess.ctx
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
func (sess *Session) event(event *WebSocketBroadcast) {
	if sess.ready {
		sess.lock.RLock()
		for subject := range sess.subscriptions {
			if subject.verify(event.Subject) {
				sess.writeEvent(MethodBroadcast, event)
				break
			}
		}
		sess.lock.RUnlock()
	}
}

// topics .
func (sess *Session) topics() []subject {
	sess.lock.RLock()

	var topics = make([]subject, 0, len(sess.subscriptions))
	for subject := range sess.subscriptions {
		topics = append(topics, subject)
	}
	sess.lock.RUnlock()

	return topics
}

// subscribe .
func (sess *Session) subscribe(params *json.RawMessage) {
	var subjects = make([]subject, 0)

	if err := json.Unmarshal(*params, &subjects); err != nil {
		logger.Context(sess.ctx).Errorf("[WebSocketID: %s] subscribe failed. %s", sess.id, err)
		sess.WriteError(webserver.StatusParamInvalid.WebError(err))
		sess.close()
		return
	}

	var subscribers = make([]*subscriber, 0, len(subjects))

	sess.lock.Lock()
	for _, subject := range subjects {
		sess.subscriptions[subject] = true

		subscribers = append(
			subscribers, &subscriber{
				sessionID: sess.id,
				subject:   subject,
				onEvent:   sess.broadcast,
			},
		)
	}
	sess.lock.Unlock()

	go es.subscribe(subscribers...)

	sess.writeEvent(MethodSubscribe, sess.topics())
}

// unsubscribe .
func (sess *Session) unsubscribe(params *json.RawMessage) {
	var subjects = make([]subject, 0)

	if err := json.Unmarshal(*params, &subjects); err != nil {
		logger.Context(sess.ctx).Errorf("[WebSocketID: %s] unsubscribe failed. %s", sess.id, err)
		sess.WriteError(webserver.StatusParamInvalid.WebError(err))
		sess.close()
		return
	}

	var subscribers = make([]*subscriber, 0, len(subjects))

	sess.lock.Lock()
	for _, subject := range subjects {
		delete(sess.subscriptions, subject)

		subscribers = append(
			subscribers, &subscriber{
				sessionID: sess.id,
				subject:   subject,
				onEvent:   sess.broadcast,
			},
		)
	}
	sess.lock.Unlock()

	go es.unsubscribe(subscribers...)

	sess.writeEvent(MethodUnsubscribe, sess.topics())
}

// close exec session.disconnect
func (sess *Session) close() {
	sess.closing <- struct{}{}
}

// disconnect .
func (sess *Session) disconnect() {
	sess.once.Do(
		func() {
			logger.Context(sess.ctx).Debugf("[WebSocketID: %s] disconnected", sess.id)

			sess.ready = false
			sess.closed = true

			store.dropSession(sess.id)

			close(sess.closing)

			close(sess.request)
			close(sess.response)
			close(sess.broadcast)

			if sess.conn != nil {
				sess.conn.Close()
				sess.conn = nil
			}

			sess = nil
		},
	)
}

var store = &pool{store: make(map[hfwctx.ID]struct{}, 0)}

// pool .
type pool struct {
	lk    sync.RWMutex
	store map[hfwctx.ID]struct{}
}

// verifySession .
func (pool *pool) verifySession(id hfwctx.ID) bool {
	pool.lk.RLock()
	_, found := pool.store[id]
	pool.lk.RUnlock()

	return found
}

// newSession .
func (pool *pool) createSession() hfwctx.ID {
	id := hfwctx.NewID()

	pool.lk.Lock()
	pool.store[id] = struct{}{}
	pool.lk.Unlock()

	return id
}

// dropSession .
func (pool *pool) dropSession(id hfwctx.ID) {
	pool.lk.Lock()
	delete(pool.store, id)
	pool.lk.Unlock()
}
