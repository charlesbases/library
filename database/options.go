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
	// DefaultContext default Options.Context
	DefaultContext = context.Background()
	// DefaultTimeout default Options.Timeout
	DefaultTimeout = time.Second * 3
	// DefaultMaxIdleConns default Options.MaxIdleConns
	DefaultMaxIdleConns = 1 << 4
	// DefaultMaxOpenConns default Options.MaxOpenConns
	DefaultMaxOpenConns = 1 << 10
	// DefaultConnMaxIdleTime default Options.ConnMaxIdleTime
	DefaultConnMaxIdleTime = 0
	// DefaultConnMaxLifetime default Options.ConnMaxLifetime
	DefaultConnMaxLifetime = 3600
)

// Options .
type Options struct {
	// Addrs addrs
	Addrs []string
	// Timeout 超时时间
	Timeout time.Duration
	// MaxIdleConns 最大空闲连接数
	MaxIdleConns int
	// MaxOpenConns 最大连接数
	MaxOpenConns int
	// ConnMaxIdleTime 连接最大空闲时间
	ConnMaxIdleTime time.Duration
	// ConnMaxLifetime 连接可复用的最大时间（秒）
	ConnMaxLifetime time.Duration
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

// SetMaxIdleConns .
func SetMaxIdleConns(n int) Option {
	return func(opts *Options) {
		opts.MaxIdleConns = n
	}
}

// SetMaxOpenConns .
func SetMaxOpenConns(n int) Option {
	return func(opts *Options) {
		opts.MaxOpenConns = n
	}
}

// SetConnMaxIdleTime .
func SetConnMaxIdleTime(d int64) Option {
	return func(opts *Options) {
		opts.ConnMaxIdleTime = time.Duration(d) * time.Second
	}
}

// SetConnMaxLifetime .
func SetConnMaxLifetime(d int64) Option {
	return func(opts *Options) {
		opts.ConnMaxLifetime = time.Duration(d) * time.Second
	}
}

// ShowSQL .
func ShowSQL(b bool) Option {
	return func(opts *Options) {
		opts.ShowSQL = b
	}
}
