package storage

import (
	"context"
	"time"
)

const (
	// defaultRegion custom
	defaultRegion = "custom"
	// defaultCallerSkip default of func caller skip
	defaultCallerSkip = 1
	// defaultTimeout 3 * time.Second
	defaultTimeout = 3 * time.Second
	// defaultListMaxKeys 默认获取所有对象
	defaultListMaxKeys = -1
	// defaultListRecursive 默认获取子文件夹对象
	defaultListRecursive = true
	// defaultPresignExpire default
	defaultPresignExpire = time.Hour * 24
)

// defaultContext default context
var defaultContext = context.Background()

// Options .
type Options struct {
	// Context context
	Context context.Context
	// Timeout 连接超时时间
	Timeout time.Duration
	// Region region
	Region string
	// SSL use ssl
	UseSSL bool
}

// NewOptions .
func NewOptions(opts ...func(o *Options)) *Options {
	var o = &Options{
		Context: defaultContext,
		Timeout: defaultTimeout,
		Region:  defaultRegion,
		UseSSL:  false,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// PutOptions .
type PutOptions struct {
	Context context.Context
	// increases the number of callers skipped by caller annotation
	CallerSkip int
}

// NewPutOptions .
func NewPutOptions(opts ...func(o *PutOptions)) *PutOptions {
	var o = &PutOptions{
		Context:    defaultContext,
		CallerSkip: defaultCallerSkip,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// GetOptions .
type GetOptions struct {
	Context context.Context
	// object version
	VersionID string
	// increases the number of callers skipped by caller annotation
	CallerSkip int
}

// NewGetOptions .
func NewGetOptions(opts ...func(o *GetOptions)) *GetOptions {
	var o = &GetOptions{
		Context:    defaultContext,
		CallerSkip: defaultCallerSkip,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// DelOptions .
type DelOptions struct {
	Context context.Context
	// object version id
	VersionID string
	// increases the number of callers skipped by caller annotation
	CallerSkip int
}

// NewDelOptions .
func NewDelOptions(opts ...func(o *DelOptions)) *DelOptions {
	var o = &DelOptions{
		Context:    defaultContext,
		CallerSkip: defaultCallerSkip,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// ListOptions .
type ListOptions struct {
	Context context.Context
	MaxKeys int
	// ignore '/' delimiter
	Recursive bool
	// increases the number of callers skipped by caller annotation
	CallerSkip int
}

// NewListOptions .
func NewListOptions(opts ...func(o *ListOptions)) *ListOptions {
	var o = &ListOptions{
		Context:    defaultContext,
		MaxKeys:    defaultListMaxKeys,
		Recursive:  defaultListRecursive,
		CallerSkip: defaultCallerSkip,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// CopyOptions .
type CopyOptions struct {
	Context context.Context
	// increases the number of callers skipped by caller annotation
	CallerSkip int
}

// NewCopyOptions .
func NewCopyOptions(opts ...func(o *CopyOptions)) *CopyOptions {
	var o = &CopyOptions{
		Context:    defaultContext,
		CallerSkip: defaultCallerSkip,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// PresignOptions .
type PresignOptions struct {
	Context context.Context
	// object version id
	VersionID string
	// expires time (s)
	Expires time.Duration
	// increases the number of callers skipped by caller annotation
	CallerSkip int
}

// NewPresignOptions .
func NewPresignOptions(opts ...func(o *PresignOptions)) *PresignOptions {
	var o = &PresignOptions{
		Context:    defaultContext,
		Expires:    defaultPresignExpire,
		CallerSkip: defaultCallerSkip,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
