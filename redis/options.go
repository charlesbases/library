package redis

import (
	"context"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/charlesbases/library/codec"
	"github.com/charlesbases/library/codec/json"
)

var (
	defaultContext   = context.Background()
	defaultTimeout   = 3 * time.Second
	defaultRetries   = 0
	defaultMarshaler = json.Marshaler

	defaultMutexTTL      = 5 * time.Second
	defaultMutexInterval = 50 * time.Millisecond

	// RandomExpiry 0-24h 随机过期时间
	RandomExpiry = func() time.Duration {
		return time.Duration(rand.Intn(24 * int(time.Hour)))
	}
)

// RedisClient redis client
var RedisClient Cmdable = func(o *Options) (redis.Cmdable, func() error) {
	client := redis.NewClient(&redis.Options{
		Addr:         o.Addrs[0],
		Username:     o.Username,
		Password:     o.Password,
		MaxRetries:   o.MaxRetries,
		DialTimeout:  o.Timeout,
		ReadTimeout:  o.Timeout,
		WriteTimeout: o.Timeout,
		// PoolSize:     4 * runtime.NumCPU(), // 连接池最大 socket 连接数。(default: 4 * runtime.NumCPU())
		// MinIdleConns: 16,                   // 空闲连接数
		// TLSConfig: &tls.Config{
		// 	InsecureSkipVerify: true,
		// },
	})

	return client, client.Close
}

// RedisCluster redis cluster
var RedisCluster Cmdable = func(o *Options) (redis.Cmdable, func() error) {
	cluster := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        o.Addrs,
		Username:     o.Username,
		Password:     o.Password,
		MaxRetries:   o.MaxRetries,
		DialTimeout:  o.Timeout,
		ReadTimeout:  o.Timeout,
		WriteTimeout: o.Timeout,
		// PoolSize:     4 * runtime.NumCPU(), // 连接池最大 socket 连接数。(default: 4 * runtime.NumCPU())
		// MinIdleConns: 16, // 空闲连接数
		// TLSConfig: &tls.Config{
		// 	InsecureSkipVerify: true,
		// },
	})

	return cluster, cluster.Close
}

// Cmdable redis.Cmdable
type Cmdable func(o *Options) (redis.Cmdable, func() error)

// Options .
type Options struct {
	Context  context.Context
	Addrs    []string
	Username string
	Password string
	// Timeout second
	Timeout time.Duration
	// MaxRetries 命令执行失败时的重试次数
	MaxRetries int

	// Cmdable redis redis.Cmdable
	Cmdable Cmdable
}

// SetOptions .
type SetOptions struct {
	Context context.Context
	// TTL key 过期时间
	TTL time.Duration
	// Expiry key 过期时间。TTL 和 Expiry 同时设置时，以 Expiry 为准
	Expiry time.Time
	// Marshaler value 编码方式
	Marshaler codec.Marshaler
}

// setoptions .
func setoptions(opts ...func(o *SetOptions)) *SetOptions {
	o := &SetOptions{
		Context:   defaultContext,
		Marshaler: defaultMarshaler,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// GetOptions .
type GetOptions struct {
	Context context.Context
	// Marshaler value 编码方式
	Marshaler codec.Marshaler
}

// getoptions .
func getoptions(opts ...func(o *GetOptions)) *GetOptions {
	o := &GetOptions{
		Context:   defaultContext,
		Marshaler: defaultMarshaler,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// DelOptions .
type DelOptions struct {
	Context context.Context
}

// deloptions .
func deloptions(opts ...func(o *DelOptions)) *DelOptions {
	o := &DelOptions{
		Context: defaultContext,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// MutexOptions .
type MutexOptions struct {
	Context context.Context
	// Interval 尝试获取锁的时间间隔
	Interval time.Duration
	// TTL 加锁后，超时自动删除锁
	TTL time.Duration
}

// mutexoptions .
func mutexoptions(opts ...func(o *MutexOptions)) *MutexOptions {
	o := &MutexOptions{
		Context:  defaultContext,
		Interval: defaultMutexInterval,
		TTL:      defaultMutexTTL,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// ExpireOptions .
type ExpireOptions struct {
	Context context.Context
	// TTL key 过期时间
	TTL time.Duration
	// Expiry key 过期时间。TTL 和 Expiry 同时设置时，以 Expiry 为准
	Expiry time.Time
}

// expireoptions .
func expireoptions(opts ...func(o *ExpireOptions)) *ExpireOptions {
	o := &ExpireOptions{
		Context: defaultContext,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
