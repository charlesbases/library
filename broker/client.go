package broker

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/charlesbases/library/codec"
	"github.com/charlesbases/library/codec/json"
)

var (
	// defaultContext default context
	defaultContext = context.Background()
	// defaultReconnectTime 重连等待时间
	defaultReconnectTime = time.Second * 3

	// ErrInvalidAddrs .
	ErrInvalidAddrs = errors.New("broker: invalid addrs")
)

// Client .
type Client interface {
	// Publish 消息发布
	// 若发布消息格式为 json, 则参数 'v' 为 JsonMessage.Data
	// 若发布消息格式为 proto, 则参数 'v' 为发布的完整消息体, 方法内不做额外的封装
	Publish(topic string, v interface{}, opts ...func(o *PublishOptions)) error
	// Subscribe 消息订阅
	Subscribe(topic string, handler Handler, otps ...func(o *SubscribeOptions))
	// Close .
	Close()
	// ID client id
	ID() string
}

type Header map[string]string

type Handler func(event Event) error

// Event .
type Event interface {
	// Topic .
	Topic() string
	// Body return bytes of message
	Body() []byte
	// Unmarshal unmarshal message
	// 若消息格式为 json, 则反序列化 message.data; 若消息格式为 proto, 则反序列化 message
	Unmarshal(v interface{}) error
}

// JsonMessage .
type JsonMessage struct {
	// ID topic 唯一标识符. uuid.NewString()
	ID string `json:"id"`
	// Producer .
	Producer string `json:"producer"`
	// CreatedAt 创建时间
	CreatedAt string `json:"created_at"`
	// Data .
	Data interface{} `json:"data"`
}

// Options .
type Options struct {
	// Address adress
	Address string
	// ReconnectTime 重连等待时间。单位：秒
	ReconnectTime time.Duration
	// Version sarama.KafkaVersion
	Version string
	// Debug print message
	Debug bool
}

// ParseOptions .
func ParseOptions(opts ...func(o *Options)) *Options {
	o := &Options{
		ReconnectTime: defaultReconnectTime,
		Debug:         false,
	}

	for _, opt := range opts {
		opt(o)
	}
	return o
}

// PublishOptions .
type PublishOptions struct {
	// Context ctx
	Context context.Context
	// Codec 序列化方式. default codec.MarshalerType_Json
	Codec codec.Marshaler
}

// ParsePublishOptions .
func ParsePublishOptions(opts ...func(o *PublishOptions)) *PublishOptions {
	o := &PublishOptions{
		Codec: json.Marshaler,
	}

	for _, opt := range opts {
		opt(o)
	}
	return o
}

// ConsumerModel 消费者模式
type ConsumerModel func(c Client, topic string) string

// RandomConsumption 随机消费。只有一个服务会消费
var RandomConsumption = func(c Client, topic string) string {
	if args := strings.Split(c.ID(), "."); len(args) != 0 {
		return args[0] + "." + topic
	}
	return topic
}

// SharedConsumption 共享消费。多服务共同消费
var SharedConsumption = func(c Client, topic string) string {
	return c.ID() + "." + topic
}

// SubscribeOptions .
type SubscribeOptions struct {
	// Context ctx
	Context context.Context
	// Codec 序列化方式. default codec.MarshalerType_Json
	Codec codec.Marshaler
	// ConsumerModel 消费者模式。多副本情况下，订阅相同 topic 的消费者是否共同处理数据
	ConsumerModel ConsumerModel
}

// ParseSubscribeOptions .
func ParseSubscribeOptions(opts ...func(o *SubscribeOptions)) *SubscribeOptions {
	o := &SubscribeOptions{
		Codec:         json.Marshaler,
		Context:       defaultContext,
		ConsumerModel: SharedConsumption,
	}

	for _, opt := range opts {
		opt(o)
	}
	return o
}
