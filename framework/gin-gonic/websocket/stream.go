package websocket

import (
	"encoding/json"
	"time"

	"github.com/charlesbases/logger"
	"github.com/gorilla/websocket"

	"github.com/charlesbases/library/framework/gin-gonic/hfwctx"
	"github.com/charlesbases/library/framework/gin-gonic/webserver"
)

const (
	defaultBuffer    = 16
	defaultTimeout   = time.Second * 3
	defaultHeartbeat = time.Second * 30
)

// Method of ws
type Method string

const (
	// MethodPing ping
	MethodPing Method = "Ping"
	// MethodResponse response
	MethodResponse Method = "Response"
	// MethodSubscribe 消息订阅
	MethodSubscribe Method = "Subscribe"
	// MethodUnsubscribe 取消订阅
	MethodUnsubscribe Method = "Unsubscribe"
	// MethodBroadcast 广播
	MethodBroadcast Method = "Broadcast"
	// MethodDisconnect 断开连接
	MethodDisconnect Method = "Disconnect"
)

// String .
func (m *Method) String() string {
	return string(*m)
}

var defaultAuth = func(c *hfwctx.Context) error { return nil }

// stream .
type stream struct {
	opts *Options
}

type (
	// WebSocketRequest .
	WebSocketRequest struct {
		ID     hfwctx.ID        `json:"id" validate:"required"`
		Method Method           `json:"method" validate:"required"`
		Params *json.RawMessage `json:"params,omitempty"`
	}
	// WebSocketResponse .
	WebSocketResponse struct {
		ID      hfwctx.ID      `json:"id" validate:"required"`
		Code    webserver.Code `json:"code,omitempty"`
		Message string         `json:"message,omitempty"`
		Method  Method         `json:"method,omitempty"`
		Data    interface{}    `json:"data,omitempty"`
	}
	// WebSocketBroadcast .
	WebSocketBroadcast struct {
		Subject subject     `json:"subject" validate:"required"`
		Time    string      `json:"time" validate:"required"`
		Data    interface{} `json:"data" validate:"required"`
	}
)

// NewStream .
func NewStream(opts ...func(o *Options)) func(c *hfwctx.Context) {
	s := &stream{opts: &Options{
		Auth:      defaultAuth,
		Buffer:    defaultBuffer,
		Timeout:   defaultTimeout,
		Heartbeat: defaultHeartbeat,
	}}
	for _, o := range opts {
		o(s.opts)
	}

	return func(c *hfwctx.Context) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		if c.Request.Header.Get("Sec-WebSocket-Protocol") != "" {
			c.Writer.Header().Set("Sec-WebSocket-Protocol", c.Request.Header.Get("Sec-WebSocket-Protocol"))
		}

		if err := s.opts.Auth(c); err != nil {
			logger.WithContext(c).Error("websocket authorization failed. ", err)
			c.ReturnError(webserver.StatusAccessDenied, err)
			return
		}

		if err := s.connect(c); err != nil {
			logger.WithContext(c).Error("websocket connect failed. ", err)
		}
		return
	}
}

// session .
func (stream *stream) newSession(c *hfwctx.Context, conn *websocket.Conn) *Session {
	return &Session{
		id:            store.createSession(),
		subscriptions: make(map[subject]bool),
		request:       make(chan *WebSocketRequest, stream.opts.Buffer),
		response:      make(chan *WebSocketResponse, stream.opts.Buffer),
		broadcast:     make(chan *WebSocketBroadcast, stream.opts.Buffer),
		ctx:           c,
		conn:          conn,
		opts:          stream.opts,
		closing:       make(chan struct{}),
	}
}

// connect .
func (stream *stream) connect(c *hfwctx.Context) error {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.WithContext(c).Error("websocket upgrade error: ", err)
		return webserver.StatusBadRequest.WebError(err)
	}
	defer conn.Close()

	sess := stream.newSession(c, conn)
	logger.WithContext(sess.ctx).Debugf("[WebSocketID: %s] connected", sess.id)

	sess.ping()
	sess.serve()
	return nil
}

// Options .
type Options struct {
	Auth      func(c *hfwctx.Context) error             // WebSocket 认证
	Buffer    int                                       // 缓冲区间。default: 16
	Timeout   time.Duration                             // 超时时间。default: 3s
	Heartbeat time.Duration                             // 心跳。default: 5s
	Action    func(c *hfwctx.Context, session *Session) // 主动推送
}
