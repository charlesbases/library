package broker

import (
	"context"
	"crypto/tls"

	"library/codec"
)

// Options .
type Options struct {
	Addrs []string

	TLSConfig *tls.Config

	// Codec 编码类型. default: "application/json"
	Codec codec.Marshaler

	Context context.Context
}

type Option func(o *Options)

// NewOptions .
func NewOptions(opts ...Option) *Options {
	options := new(Options)
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// WithAddrs .
func WithAddrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// WithTLSConfig .
func WithTLSConfig(tls *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = tls
	}
}

// WithCodec .
func WithCodec(codec codec.Marshaler) Option {
	return func(o *Options) {
		o.Codec = codec
	}
}

// WithContext .
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// PublishOptions 消息推送
type PublishOptions struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type PublishOption func(o *PublishOptions)

// NewPublishOptions .
func NewPublishOptions(opts ...PublishOption) *PublishOptions {
	options := new(PublishOptions)
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// WithPublishContext .
func WithPublishContext(ctx context.Context) PublishOption {
	return func(o *PublishOptions) {
		o.Context = ctx
	}
}

// SubscribeOptions 消息订阅
type SubscribeOptions struct {
	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is acked.
	AutoAck bool
	// Subscribers with the same queue name
	// will create a shared subscription where each
	// receives a subset of messages.
	Queue string

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type SubscribeOption func(o *SubscribeOptions)

// NewSubscribeOptions .
func NewSubscribeOptions(opts ...SubscribeOption) *SubscribeOptions {
	options := &SubscribeOptions{
		AutoAck: true,
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// WithSubscribeQueue .
func WithSubscribeQueue(queue string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = queue
	}
}

// WithSubscribeContext .
func WithSubscribeContext(ctx context.Context) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Context = ctx
	}
}
