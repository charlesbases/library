package broker

import (
	"errors"

	"library/codec"

	"github.com/google/uuid"
)

var (
	// ErrInvalidMsg .
	ErrInvalidMsg = errors.New("broker: invalid message or message nil")
	// ErrInvalidAddrs .
	ErrInvalidAddrs = errors.New("broker: invalid addrs")
)

type Broker interface {
	// Address .
	Address() []string
	// Publish 消息发布
	Publish(topic string, v interface{}, opts ...PublishOption) error
	// Subscribe 消息订阅
	Subscribe(topic string, handler Handler, otps ...SubscribeOption)
	// OnStart lifecycle.Hook
	OnStart() error
	// OnStop lifecycle.Hook
	OnStop() error
	// String .
	String() string
}

type Header map[string]string

type Handler func(event Event) error

// Message .
type Message struct {
	// ID topic 唯一标识符 16字节
	ID uuid.UUID `json:"id"`
	// Topic subject
	Topic string `json:"topic"`
	// Producer .
	Producer string `json:"producer"`
	// CreatedAt 创建时间
	CreatedAt string `json:"created_at"`
	// ContentType 编码类型 application/json | application/proto
	ContentType codec.ContentType `json:"content_type"`
	// Data .
	Data interface{} `json:"data"`
}

// Event .
type Event interface {
	// Topic .
	Topic() string
	// Body return bytes of Message.Data
	Body() []byte
	// Unmarshal unmarshal for Message.Data. must be a pointer
	Unmarshal(v interface{}) error
}
