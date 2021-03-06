package websocket

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	logger "library/logger/seelog"
	"library/websocket/pb"

	"github.com/google/uuid"
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
		ID     string           `json:"id" validate:"required"`
		Method pb.Method        `json:"method" validate:"required"`
		Params *json.RawMessage `json:"params,omitempty"`
	}
	// WebSocketResponse .
	WebSocketResponse struct {
		ID      string      `json:"id" validate:"required"`
		Code    int         `json:"code" validate:"required"`
		Message string      `json:"message,omitempty"`
		Method  pb.Method   `json:"method,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}
	// WebSocketBroadcast .
	WebSocketBroadcast struct {
		Topic topic       `json:"topic" validate:"required"`
		Data  interface{} `json:"data" validate:"required"`
	}
)

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

// NewStream .
func NewStream(opts ...option) *stream {
	stream := new(stream)
	stream.init(opts...)
	return stream
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

// options .
type options struct {
	auth      func(r *http.Request) error // ??????
	buffer    int                         // ???????????????default: 16
	timeout   time.Duration               // ???????????????default: 3s
	heartbeat time.Duration               // ?????????default: 5s
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

var store = &pool{store: make(map[string]struct{}, 0)}

// pool .
type pool struct {
	lk    sync.RWMutex
	store map[string]struct{}
}

// verifySession .
func (pool *pool) verifySession(id string) bool {
	var active bool
	pool.lk.RLock()
	_, active = pool.store[id]
	pool.lk.RUnlock()
	return active
}

// newSession .
func (pool *pool) createSession() string {
	id := uuid.New().String()

	pool.lk.Lock()
	pool.store[id] = struct{}{}
	pool.lk.Unlock()

	return id
}

// dropSession .
func (pool *pool) dropSession(id string) {
	pool.lk.Lock()
	delete(pool.store, id)
	pool.lk.Unlock()
}
