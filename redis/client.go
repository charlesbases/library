package redis

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/charlesbases/logger"

	"github.com/charlesbases/library"
)

// ErrRedisNotReady redis is not ready
var ErrRedisNotReady = errors.New("redis is not ready")

var (
	// lockPrefixKey redis 分布式锁的 key 前缀
	lockPrefixKey = KeyPrefix("lock_")
)

type keyword string

// Key keyword of key
var Key = func(key string) keyword {
	return keyword(key)
}

// KeyPrefix keyword of prefix
var KeyPrefix = func(prefix string) func(key string) keyword {
	return func(key string) keyword {
		var builder strings.Builder
		builder.WriteString(prefix)
		builder.WriteString(key)
		return keyword(builder.String())
	}
}

// KeySuffix keyword of suffix
var KeySuffix = func(suffix string) func(key string) keyword {
	return func(key string) keyword {
		var builder strings.Builder
		builder.WriteString(key)
		builder.WriteString(suffix)
		return keyword(builder.String())
	}
}

// Mutex redis 分布式锁
type Mutex struct {
	key  keyword
	opts *MutexOptions

	err    error
	locked bool
}

// Err .
func (m *Mutex) Err() error {
	return m.err
}

// Lock .
func (m *Mutex) Lock() {
	t := time.NewTicker(m.opts.Interval)

	for {
		select {
		case <-t.C:
			ok, _ := r.client.SetNX(m.opts.Context, string(m.key), r.id, m.opts.TTL).Result()
			if ok {
				m.locked = true
				logger.WithContext(m.opts.Context).Debugf(`[redis](%s): locked: %v.`, m.key, m.opts.TTL)
				return
			}
		}
	}
}

// Unlock .
func (m *Mutex) Unlock() {
	if !m.locked {
		logger.WithContext(m.opts.Context).Errorf(`[redis](%s): unlocked: unlock of unlocked mutex`, m.key)
		return
	}

	m.locked = false

	logger.WithContext(m.opts.Context).Debugf(`[redis](%s): unlocked`, m.key)
	r.Del(m.key, func(o *DelOptions) { o.Context = m.opts.Context })
}

// r default client
var r *rkv

// rkv .
type rkv struct {
	id      string
	opts    *Options
	client  redis.Cmdable
	active  bool
	closing func() error
}

// isReady .
func (r *rkv) isReady() bool {
	if r != nil && r.active {
		return true
	}
	return false
}

// RedisMessage .
type RedisMessage struct {
	Data      interface{} `json:"data"`
	CreatedBy string      `json:"created_by"`
	CreatedAt string      `json:"created_at"`
}

// JsonWrap 将 val 包装成 RedisMessage
func JsonWrap(val interface{}) *RedisMessage {
	return &RedisMessage{
		Data:      val,
		CreatedBy: r.id,
		CreatedAt: library.NowString(),
	}
}

// Set .
func (r *rkv) Set(key keyword, val interface{}, opts ...func(o *SetOptions)) *StatusOutput {
	var sopts = setoptions(opts...)

	output := &StatusOutput{baseOutput: baseOutput{ctx: sopts.Context, key: string(key)}}
	if !r.isReady() {
		output.err = errors.Errorf("[redis](%s): set: %v", key, ErrRedisNotReady)
		return output
	}

	if !sopts.Expiry.IsZero() {
		sopts.TTL = time.Since(sopts.Expiry)
	}

	data, err := sopts.Marshaler.Marshal(val)
	if err != nil {
		output.err = errors.Errorf("[redis](%s): set: %v", key, err)
		return output
	}

	if err := r.client.Set(sopts.Context, string(key), data, sopts.TTL).Err(); err != nil {
		output.err = errors.Errorf("[redis](%s): set: %v", key, err)
		return output
	}

	return output
}

// Get .
func (r *rkv) Get(key keyword, opts ...func(o *GetOptions)) *BytesOutput {
	var gopts = getoptions(opts...)

	output := &BytesOutput{marshaler: gopts.Marshaler, baseOutput: baseOutput{ctx: gopts.Context, key: string(key)}}
	if !r.isReady() {
		output.err = errors.Errorf("[redis](%s): get: %v", key, ErrRedisNotReady)
		return output
	}

	data, err := r.client.Get(gopts.Context, string(key)).Bytes()
	if err != nil {
		output.err = errors.Errorf("[redis](%s): get: %v", key, err)
		return output
	}

	ttl, err := r.client.TTL(gopts.Context, string(key)).Result()
	if err != nil {
		output.err = errors.Errorf("[redis](%s): get.ttl: %v", key, err)
		return output
	}

	output.val = data
	output.ttl = ttl
	output.expiry = time.Now().Add(ttl)
	return output
}

// Del .
func (r *rkv) Del(key keyword, opts ...func(o *DelOptions)) *StatusOutput {
	var dopts = deloptions(opts...)

	output := &StatusOutput{baseOutput: baseOutput{ctx: dopts.Context, key: string(key)}}
	if !r.isReady() {
		output.err = errors.Errorf("[redis](%s): del: %v", key, ErrRedisNotReady)
		return output
	}

	// 调用 redis.Del() 进行删除
	// 注意：redis 的删除策略为惰性删除，并不确保立即删除，并且删除键值对会占用 CPU 资源，尤其是大量删除时
	// if err := r.client.Del(dopts.Context, string(key)).Err(); err != nil {
	// 	output.err = errors.Errorf("[redis](%s): del: %v", key, err)
	// 	return output
	// }

	// 使用 redis.RenameNX() 后设置过期时间的方式进行平滑删除
	// 相较于 redis.Del()，定时删除可以在一定程度上分摊删除操作的 CPU 负载
	_, err := r.client.TxPipelined(dopts.Context, func(pipe redis.Pipeliner) error {
		newkey := uuid.NewString()
		if err := pipe.RenameNX(dopts.Context, string(key), newkey).Err(); err != nil {
			return err
		}
		// 将 key 的过期时间(删除时间)设为 0-3s, 防止集中删除
		return pipe.PExpire(dopts.Context, newkey, time.Duration(rand.Intn(3000))*time.Millisecond).Err()
	})

	output.err = errors.Wrapf(err, "[redis](%s): del", key)
	return output
}

// Expire .
func (r *rkv) Expire(key keyword, opts ...func(o *ExpireOptions)) *StatusOutput {
	var eopts = expireoptions(opts...)

	output := &StatusOutput{baseOutput: baseOutput{ctx: eopts.Context, key: string(key)}}
	if !r.isReady() {
		output.err = errors.Errorf("[redis](%s): expire: %v", key, ErrRedisNotReady)
		return output
	}

	if !eopts.Expiry.IsZero() {
		// ExpireAt
		output.err = r.client.PExpireAt(eopts.Context, string(key), eopts.Expiry).Err()
	} else {
		// TTL
		output.err = r.client.PExpire(eopts.Context, string(key), eopts.TTL).Err()
	}

	output.err = errors.Wrapf(output.err, "[redis](%s): expire", key)
	return output
}

// Mutex .
func (r *rkv) Mutex(key keyword, opts ...func(o *MutexOptions)) *Mutex {
	if !r.isReady() {
		return &Mutex{err: errors.Errorf("[redis](%s): mutex: %v", key, ErrRedisNotReady)}
	}

	return &Mutex{
		key:  lockPrefixKey(string(key)),
		opts: mutexoptions(opts...),
	}
}

// IsExists .
func (r *rkv) IsExists(key string, opts ...func(o *GetOptions)) *BoolOutput {
	var gopts = getoptions(opts...)

	output := &BoolOutput{baseOutput: baseOutput{ctx: gopts.Context, key: string(key)}}
	if !r.isReady() {
		output.err = errors.Errorf("[redis](%s): exists: %v", key, ErrRedisNotReady)
		return output
	}

	if r.client.Exists(gopts.Context, key).Val() != 0 {
		output.val = true
	}
	return output
}

// Close .
func (r *rkv) Close() error {
	if r.closing != nil {
		return r.closing()
	}
	return nil
}

// ping .
func (r *rkv) ping() error {
	if err := r.client.Ping(r.opts.Context).Err(); err != nil {
		r.active = false
		return errors.Wrapf(err, "[redis]: ping")
	}
	r.active = true
	return nil
}

// nodisplay .
type nodisplay struct{}

// Printf .
func (l *nodisplay) Printf(_ context.Context, _ string, _ ...interface{}) {}

// NewClient .
func NewClient(id string, opts ...func(o *Options)) (*rkv, error) {
	redis.SetLogger(new(nodisplay))

	options := &Options{
		Cmdable:    RedisClient,
		Addrs:      []string{"0.0.0.0:6379"},
		Context:    defaultContext,
		Timeout:    defaultTimeout,
		MaxRetries: defaultRetries,
	}
	for _, opt := range opts {
		opt(options)
	}

	client, closing := options.Cmdable(options)

	r := &rkv{
		id:      id,
		opts:    options,
		client:  client,
		closing: closing,
	}

	if err := r.ping(); err != nil {
		return nil, err
	}
	return r, nil
}

// Init init default redis
func Init(id string, opts ...func(o *Options)) error {
	if client, err := NewClient(id, opts...); err != nil {
		return err
	} else {
		r = client
	}
	return nil
}

// Client .
func Client() *rkv {
	return r
}

// Close .
func Close() error {
	if r != nil {
		return r.Close()
	}
	return nil
}
