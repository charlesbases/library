package database

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrorInvaildAddrs invalid addrs
	ErrorInvaildAddrs = errors.New("invalid addrs")
	// ErrorDatabaseNil db is not initialized or closed
	ErrorDatabaseNil = errors.New("db is not initialized or closed")
)

var (
	// DefaultContext default context
	DefaultContext = context.Background()
	// DefaultTimeout default timeout
	DefaultTimeout = time.Second * 3
)

// Options .
type Options struct {
	// Addrs addrs
	Addrs []string
	// Timeout 超时时间
	Timeout time.Duration
	// Context context
	Context context.Context
	// ShowSQL 是否显示 sql 日志
	ShowSQL bool
}

// DefaultOptions .
func DefaultOptions() *Options {
	return &Options{
		Context: DefaultContext,
		ShowSQL: false,
	}
}

type Option func(opts *Options)

// Addrs .
func Addrs(addrs ...string) Option {
	return func(opts *Options) {
		opts.Addrs = addrs
	}
}

// Timeout .
func Timeout(t time.Duration) Option {
	return func(opts *Options) {
		opts.Timeout = t
	}
}

// ShowSQL .
func ShowSQL(b bool) Option {
	return func(opts *Options) {
		opts.ShowSQL = b
	}
}
