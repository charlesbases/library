package broker

import (
	"context"
	"time"

	"library/codec"
	"library/codec/json"
	"library/codec/proto"
)

var (
	// DefaultContext default context
	DefaultContext = context.Background()
	// DefaultReconnectTime 重连等待时间
	DefaultReconnectTime = time.Second * 3
)

var (
	// CodecJson json 编码
	CodecJson = json.NewMarshaler()
	// CodecProto proto 编码
	CodecProto = proto.NewMarshaler()
)

// Options .
type Options struct {
	Addrs []string
	// ReconnectTime 重连等待时间。单位：秒
	ReconnectTime time.Duration
}

type Option func(o *Options)

// Addrs .
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// ReconnectTime 重连等待时间。单位：秒
func ReconnectTime(d int) Option {
	return func(o *Options) {
		o.ReconnectTime = time.Second * time.Duration(d)
	}
}

// PublishOptions .
type PublishOptions struct {
	// Codec 序列化方式. default codec.MarshalerType_Json
	Codec codec.Marshaler
}

type PublishOption func(o *PublishOptions)

// DefaultPublishOptions .
func DefaultPublishOptions() *PublishOptions {
	return &PublishOptions{
		Codec: CodecJson,
	}
}

// PublishJson .
func PublishJson() PublishOption {
	return func(o *PublishOptions) {
		o.Codec = CodecJson
	}
}

// PublishProto .
func PublishProto() PublishOption {
	return func(o *PublishOptions) {
		o.Codec = CodecProto
	}
}

// SubscribeOptions .
type SubscribeOptions struct {
	// Context ctx
	Context context.Context
	// Codec 序列化方式. default codec.MarshalerType_Json
	Codec codec.Marshaler
}

type SubscribeOption func(o *SubscribeOptions)

// DefaultSubscribeOptions .
func DefaultSubscribeOptions() *SubscribeOptions {
	return &SubscribeOptions{
		Codec:   CodecJson,
		Context: DefaultContext,
	}
}

// SubscribeJson .
func SubscribeJson() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Codec = CodecJson
	}
}

// SubscribeProto .
func SubscribeProto() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Codec = CodecProto
	}
}

// SubscribeContext .
func SubscribeContext(c context.Context) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Context = c
	}
}
