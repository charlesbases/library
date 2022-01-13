package registry

import (
	"context"
	"crypto/tls"
	"time"
)

type (
	Options struct {
		Addrs     []string
		Timeout   time.Duration
		Secure    bool
		TLSConfig *tls.Config
	}

	ListOptions struct {
		Context context.Context
	}

	RegisterOptions struct {
		TTL     time.Duration
		Context context.Context
	}

	DeregisterOptions struct {
		Context context.Context
	}

	Option           func(o *Options)
	ListOption       func(o *ListOptions)
	RegisterOption   func(o *RegisterOptions)
	DeregisterOption func(o *DeregisterOptions)
)

func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

func Secure() Option {
	return func(o *Options) {
		o.Secure = true
	}
}

func TLSConfig(tls *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = tls
	}
}

func ListContext(ctx context.Context) ListOption {
	return func(o *ListOptions) {
		o.Context = ctx
	}
}

func RegisterTTL(ttl time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.TTL = ttl
	}
}

func RegisterContext(ctx context.Context) RegisterOption {
	return func(o *RegisterOptions) {
		o.Context = ctx
	}
}

func DeregisterContext(ctx context.Context) DeregisterOption {
	return func(o *DeregisterOptions) {
		o.Context = ctx
	}
}
