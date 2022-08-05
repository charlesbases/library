package websocket

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"library/websocket/pb"

	"github.com/charlesbases/logger"
	"github.com/gorilla/websocket"
)

const (
	defaultBuffer    = 16
	defaultTimeout   = time.Second * 3
	defaultHeartbeat = time.Second * 30
)

var defaultAuth = func(r *http.Request) error { return nil }

// stream .
type stream struct {
	opts *options
}

type (
	// WebSocketRequest .
	WebSocketRequest struct {
		ID     sessionID        `json:"id" validate:"required"`
		Method pb.Method        `json:"method" validate:"required"`
		Params *json.RawMessage `json:"params,omitempty"`
	}
	// WebSocketResponse .
	WebSocketResponse struct {
		ID      sessionID   `json:"id" validate:"required"`
		Code    int         `json:"code" validate:"required"`
		Message string      `json:"message,omitempty"`
		Method  pb.Method   `json:"method,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}
	// WebSocketBroadcast .
	WebSocketBroadcast struct {
		Subject subject     `json:"subject" validate:"required"`
		Time    string      `json:"time" validate:"required"`
		Data    interface{} `json:"data" validate:"required"`
	}
)

// WebError .
type WebError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Error .
func (err *WebError) Error() string {
	data, _ := json.Marshal(err)
	return string(data)
}

// NewStream .
func NewStream(opts ...option) *stream {
	stream := new(stream)
	stream.init(opts...)
	return stream
}

// init .
func (s *stream) init(opts ...option) {
	var options = new(options)
	for _, o := range opts {
		o(options)
	}
	s.opts = options

	if s.opts.auth == nil {
		s.opts.auth = defaultAuth
	}
	if s.opts.buffer == 0 {
		s.opts.buffer = defaultBuffer
	}
	if s.opts.timeout == 0 {
		s.opts.timeout = defaultTimeout
	}
	if s.opts.heartbeat == 0 {
		s.opts.heartbeat = defaultHeartbeat
	}
}

func (s *stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Header.Get("Sec-WebSocket-Protocol") != "" {
		w.Header().Set("Sec-WebSocket-Protocol", r.Header.Get("Sec-WebSocket-Protocol"))
	}

	if err := s.opts.auth(r); err != nil {
		logger.Error("websocket authorization failed. ", err)
		return
	}

	if err := s.connect(w, r); err != nil {
		logger.Error("websocket connect failed. ", err)
	}
	return
}

// connect .
func (stream *stream) connect(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("websocket upgrade error: ", err)
		return &WebError{Code: http.StatusBadRequest, Message: err.Error()}
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
		subscriptions: make(map[subject]bool),
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

// options .
type options struct {
	auth      func(r *http.Request) error // 认证
	buffer    int                         // 缓冲区间。default: 16
	timeout   time.Duration               // 超时时间。default: 3s
	heartbeat time.Duration               // 心跳。default: 5s
}

type option func(o *options)

// WithAuth .
func WithAuth(auth func(r *http.Request) error) option {
	return func(o *options) {
		o.auth = auth
	}
}

// WithBuffer .
func WithBuffer(buffer int) option {
	return func(o *options) {
		o.buffer = buffer
	}
}

// WithTimeout .
func WithTimeout(d time.Duration) option {
	return func(o *options) {
		o.timeout = d
	}
}

// WithHeartbeat .
func WithHeartbeat(d time.Duration) option {
	return func(o *options) {
		o.heartbeat = d
	}
}
