package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"library/websocket/pb"

	"github.com/charlesbases/logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// defaultDateFormat default Data Format
const defaultDateFormat = "2006-01-02 15:04:05.000"

// upgrader websocker upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: time.Second * 3,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

type metadata map[string]string

// session .
type session struct {
	id            ID
	header        metadata
	subscriptions map[topic]bool

	request   chan *WebSocketRequest   // 请求
	response  chan *WebSocketResponse  // 响应
	broadcast chan *WebSocketBroadcast // 广播

	ctx  context.Context
	conn *websocket.Conn

	opts *options

	lock sync.RWMutex
	once sync.Once

	ready   bool
	closed  bool
	closing chan struct{}
}

// connect .
func (stream *stream) connect(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("websocket upgrade error: ", err)
		return &WebError{ErrCode: http.StatusBadRequest, ErrMsg: err.Error()}
	}
	defer conn.Close()

	session := stream.session(r, conn)
	logger.Debugf("[WebSocketID: %s] connected", session.id)

	session.ping()
	session.serve()
	return nil
}

// session .
func (stream *stream) session(r *http.Request, conn *websocket.Conn) *session {
	return &session{
		id:            store.createSession(),
		header:        stream.parseHeader(r),
		subscriptions: make(map[topic]bool),
		request:       make(chan *WebSocketRequest, stream.opts.buffer),
		response:      make(chan *WebSocketResponse, stream.opts.buffer),
		broadcast:     make(chan *WebSocketBroadcast, stream.opts.buffer),
		ctx:           r.Context(),
		conn:          conn,
		opts:          stream.opts,
		closing:       make(chan struct{}),
	}
}

// parseHeader .
func (stream *stream) parseHeader(r *http.Request) metadata {
	var header metadata = make(map[string]string, 0)
	for key, val := range r.Header {
		header[key] = strings.Join(val, ",")
	}
	return header
}

// isCloseError .
func (session *session) isCloseError(err error) int {
	if e, ok := err.(*websocket.CloseError); ok {
		return e.Code
	} else {
		return -1
	}
}

// serve .
func (session *session) serve() {
	go session.listening()

	for {
		select {
		case <-time.After(session.opts.heartbeat):
			if session.ready {
				session.ready = false
				logger.Debugf("[WebSocketID: %s] no heartbeat", session.id)
			}
		case request, ok := <-session.request:
			if ok {
				switch request.Method {
				case pb.Method_PING:
					session.ping()
				case pb.Method_SUBSCRIBE:
					session.subscribe(request.Params)
				case pb.Method_UNSUBSCRIBE:
					session.unsubscribe(request.Params)
				case pb.Method_DISCONNECT:
					session.close()
				default:
					session.error(http.StatusBadRequest, fmt.Sprintf("invalid method: %d", request.Method))
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
func (session *session) listening() {
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
					logger.Errorf("[WebSocketID: %s] read message error: %v", session.id, err)
					session.error(http.StatusBadRequest, err.Error())
				}

				session.close()
				return
			}

			if !store.verifySession(request.ID) {
				session.error(http.StatusBadRequest, fmt.Sprintf("invalid session id, %s not connected", request.ID))
				session.close()
				return
			}

			if session.closed {
				session.error(http.StatusBadRequest, "connect closed")
				return
			}

			if request.Method == pb.Method_PING {
				goto r
			}

			if !session.ready {
				session.error(http.StatusBadRequest, "connect not ready")
				continue
			}

			if request.Params == nil {
				session.error(http.StatusBadRequest, "params cannot be empty.")
				continue
			}

		r:
			logger.Debugf("[WebSocketID: %s] [r] Method: %s", session.id, request.Method.String())
			session.request <- request
		}
	}
}

// write .
func (session *session) write(v *WebSocketResponse) error {
	if session.ready {
		v.ID = session.id
		logger.Debugf("[WebSocketID: %s] [w] Method: %s", session.id, v.Method.String())
		return session.conn.WriteJSON(v)
	}

	logger.Debugf("[WebSocketID: %s] write failed. connect not ready.", session.id)
	return nil
}

// ping .
func (session *session) ping() {
	if !session.closed {
		session.ready = true

		session.normal(pb.Method_PING, time.Now().Format(defaultDateFormat))
	}
}

// normal WebSocket 正常返回
func (session *session) normal(method pb.Method, data interface{}) {
	session.response <- &WebSocketResponse{
		ID:      session.id,
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Method:  method,
		Data:    data,
	}
}

// error WebSocket 返回错误
func (session *session) error(errCode int, errMsg string) {
	session.response <- &WebSocketResponse{
		ID:      session.id,
		Code:    errCode,
		Message: errMsg,
	}
}

type topic string

// verify .
func (sub *topic) verify(event topic) bool {
	if *sub == "*" {
		return true
	}
	if *sub == event {
		return true
	}
	if strings.HasSuffix(string(*sub), "*") {
		prefix := strings.TrimSuffix(string(*sub), "*")
		return strings.HasPrefix(string(event), prefix)
	}
	return false
}

// event .
func (session *session) event(event *WebSocketBroadcast) {
	if session.ready {
		session.lock.RLock()
		for topic := range session.subscriptions {
			if topic.verify(event.Topic) {
				session.normal(pb.Method_BROADCAST, event)
				break
			}
		}
		session.lock.RUnlock()
	}
}

// topics .
func (session *session) topics() []topic {
	session.lock.RLock()

	var topics = make([]topic, 0, len(session.subscriptions))
	for topic := range session.subscriptions {
		topics = append(topics, topic)
	}
	session.lock.RUnlock()

	return topics
}

// subscribe .
func (session *session) subscribe(params *json.RawMessage) {
	var topics = make([]topic, 0)

	if err := json.Unmarshal(*params, &topics); err != nil {
		logger.Errorf("[WebSocketID: %s] subscribe failed. %v", session.id, err)
		session.error(http.StatusBadRequest, err.Error())
		session.close()
		return
	}

	session.lock.Lock()
	for _, topic := range topics {
		session.subscriptions[topic] = true
	}
	session.lock.Unlock()

	session.normal(pb.Method_SUBSCRIBE, session.topics())
}

// unsubscribe .
func (session *session) unsubscribe(params *json.RawMessage) {
	var topics = make([]topic, 0)

	if err := json.Unmarshal(*params, &topics); err != nil {
		logger.Errorf("[WebSocketID: %s] unsubscribe failed. %v", session.id, err)
		session.error(http.StatusBadRequest, err.Error())
		session.close()
		return
	}

	session.lock.Lock()
	for _, topic := range topics {
		delete(session.subscriptions, topic)
	}
	session.lock.Unlock()

	session.normal(pb.Method_UNSUBSCRIBE, session.topics())

}

// close exec session.disconnect
func (session *session) close() {
	session.closing <- struct{}{}
}

// disconnect .
func (session *session) disconnect() {
	session.once.Do(func() {
		logger.Debugf("[WebSocketID: %s] disconnected", session.id)

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

type ID string

var store = &pool{store: make(map[ID]struct{}, 0)}

// pool .
type pool struct {
	lk    sync.RWMutex
	store map[ID]struct{}
}

// verifySession .
func (pool *pool) verifySession(id ID) bool {
	var active bool
	pool.lk.RLock()
	_, active = pool.store[id]
	pool.lk.RUnlock()
	return active
}

// newSession .
func (pool *pool) createSession() ID {
	id := ID(uuid.New().String())

	pool.lk.Lock()
	pool.store[id] = struct{}{}
	pool.lk.Unlock()

	return id
}

// dropSession .
func (pool *pool) dropSession(id ID) {
	pool.lk.Lock()
	delete(pool.store, id)
	pool.lk.Unlock()
}
