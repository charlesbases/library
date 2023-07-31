package storage

import (
	"context"
	"time"
)

const (
	// defaultRegion self-built
	defaultRegion = "self-built"
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
	// Context .
	Context context.Context
}

// NewPutOptions .
func NewPutOptions(opts ...func(o *PutOptions)) *PutOptions {
	var o = &PutOptions{
		Context: defaultContext,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// GetOptions .
type GetOptions struct {
	// Context .
	Context context.Context
	// VersionID object version
	VersionID string
}

// NewGetOptions .
func NewGetOptions(opts ...func(o *GetOptions)) *GetOptions {
	var o = &GetOptions{
		Context: defaultContext,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// DelOptions .
type DelOptions struct {
	// Context .
	Context context.Context
	// VersionID version id
	VersionID string
}

// NewDelOptions .
func NewDelOptions(opts ...func(o *DelOptions)) *DelOptions {
	var o = &DelOptions{
		Context: defaultContext,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// ListOptions .
type ListOptions struct {
	// Context .
	Context context.Context
	// MaxKeys .
	MaxKeys int
	// Recursive Ignore '/' delimiter
	Recursive bool
}

// NewListOptions .
func NewListOptions(opts ...func(o *ListOptions)) *ListOptions {
	var o = &ListOptions{
		Context:   defaultContext,
		MaxKeys:   defaultListMaxKeys,
		Recursive: defaultListRecursive,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// CopyOptions .
type CopyOptions struct {
	// Context .
	Context context.Context
}

// NewCopyOptions .
func NewCopyOptions(opts ...func(o *CopyOptions)) *CopyOptions {
	var o = &CopyOptions{
		Context: defaultContext,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// PresignOptions .
type PresignOptions struct {
	// Context .
	Context context.Context
	// VersionID data version
	VersionID string
	// Expires expires time (s)
	Expires time.Duration
}

// NewPresignOptions .
func NewPresignOptions(opts ...func(o *PresignOptions)) *PresignOptions {
	var o = &PresignOptions{
		Context: defaultContext,
		Expires: defaultPresignExpire,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
