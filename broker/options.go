package broker

import (
	"context"
	"strings"
	"time"

	"github.com/charlesbases/library/codec"
	"github.com/charlesbases/library/codec/json"
)

// Options .
type Options struct {
	// Address adress
	Address string
	// ReconnectWait 重连等待时间。单位：秒
	ReconnectWait time.Duration
	// Version sarama.KafkaVersion
	Version string
}

// ParseOptions .
func ParseOptions(opts ...func(o *Options)) *Options {
	o := &Options{
		ReconnectWait: defaultReconnectWait,
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
type ConsumerModel func(clientid string, topic string) string

// RandomConsumption 随机消费。只有一个服务会消费
var RandomConsumption = func(clientid string, topic string) string {
	if args := strings.Split(clientid, "."); len(args) != 0 {
		return args[0] + "." + topic
	}
	return topic
}

// SharedConsumption 共享消费。多服务共同消费
var SharedConsumption = func(clientid string, topic string) string {
	return clientid + "." + topic
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
