package nats

import (
	"context"
	"crypto/tls"
	"strings"

	"library/codec"
	"library/codec/json"

	"github.com/nats-io/nats.go"
)

// Options .
type Options struct {
	Addrs   []string
	Stream  string
	Context context.Context

	// Codec 编码类型. default: "application/json"
	Codec codec.Marshaler

	TLSConfig *tls.Config
}

type Option func(o *Options)

// defaultOptions .
func defaultOptions() *Options {
	return &Options{
		Addrs:   []string{nats.DefaultURL},
		Stream:  defaultStreamName,
		Codec:   json.NewMarshaler(),
		Context: context.Background(),
	}
}

// WithAddrs .
func WithAddrs(addrs ...string) Option {
	return func(o *Options) {
		for _, addr := range addrs {
			if strings.HasPrefix(addr, "nats://") {
				o.Addrs = append(o.Addrs, addr)
			} else {
				o.Addrs = append(o.Addrs, "nats://"+addr)
			}
		}
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

// defaultPublishOptions .
func defaultPublishOptions() *PublishOptions {
	return &PublishOptions{
		Context: context.Background(),
	}
}

// WithPublishContext .
func WithPublishContext(ctx context.Context) PublishOption {
	return func(o *PublishOptions) {
		o.Context = ctx
	}
}

// SubscribeOptions 消息订阅
type SubscribeOptions struct {
	// Queue Subscribers with the same queue name
	// will create a shared subscription where each
	// receives a subset of messages.
	Queue string

	// Context Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type SubscribeOption func(o *SubscribeOptions)

// defaultSubscribeOptions .
func defaultSubscribeOptions() *SubscribeOptions {
	return &SubscribeOptions{
		Context: context.Background(),
	}
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
