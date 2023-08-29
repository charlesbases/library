package redis

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/charlesbases/library/logger"
	"github.com/charlesbases/library/once"
	"github.com/charlesbases/library/sonyflake"

	"github.com/charlesbases/library"
)

func init() {
	once.RandSeed(time.Now().UnixNano())
}

// ErrRedisNotReady redis is not ready
var ErrRedisNotReady = errors.New("redis is not ready")

var (
	// lockPrefixKey redis 分布式锁的 key 前缀
	lockPrefixKey = KeyPrefix("lock_")
	// delSuffixKey .
	delSuffixKey = func(key keyword) keyword {
		var builder strings.Builder
		builder.WriteString(string(key))
		builder.WriteString("_d")
		builder.WriteString(sonyflake.NextID().String())
		return keyword(builder.String())
	}
)

type keyword string

var Key = func(key string) keyword {
	return keyword(key)
}

var KeyPrefix = func(prefix string) func(key string) keyword {
	return func(key string) keyword {
		var builder strings.Builder
		builder.WriteString(prefix)
		builder.WriteString(key)
		return keyword(builder.String())
	}
}

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
	*rkv
	key  keyword
	err  error
	opts *MutexOptions

	tk *time.Ticker // 尝试获取锁的时间间隔
	tm *time.Timer  // 超时后自动删除
}

// Err .
func (m *Mutex) Err() error {
	return m.err
}

// Lock .
func (m *Mutex) Lock() {
	for {
		select {
		case <-m.tk.C:
			ok, _ := m.client.SetNX(m.opts.Context, string(m.key), m.opts.Mark, m.opts.TTL).Result()
			if ok {
				logger.DebugfWithContext(m.opts.Context, `[redis](%s) locked %v.`, m.key, m.opts.TTL)
				return
			}
		case <-m.tm.C:
			m.Unlock()
		}
	}
}

// Unlock .
func (m *Mutex) Unlock() {
	logger.DebugfWithContext(m.opts.Context, `[redis](%s) unlocked.`, m.key)

	m.Del(m.key, func(o *DelOptions) { o.Context = m.opts.Context })
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
		output.err = ErrRedisNotReady
		return output
	}

	if !sopts.Expiry.IsZero() {
		sopts.TTL = time.Since(sopts.Expiry)
	}

	data, err := sopts.Marshaler.Marshal(val)
	if err != nil {
		logger.ErrorfWithContext(sopts.Context, `[redis](%s) set failed. %s`, key, err.Error())
		output.err = err
		return output
	}

	if err := r.client.Set(sopts.Context, string(key), data, sopts.TTL).Err(); err != nil {
		logger.ErrorfWithContext(sopts.Context, `[redis](%s) set failed. %s`, key, err.Error())
		output.err = err
		return output
	}

	return output
}

// Get .
func (r *rkv) Get(key keyword, opts ...func(o *GetOptions)) *BytesOutput {
	var gopts = getoptions(opts...)

	output := &BytesOutput{marshaler: gopts.Marshaler, baseOutput: baseOutput{ctx: gopts.Context, key: string(key)}}
	if !r.isReady() {
		output.err = ErrRedisNotReady
		return output
	}

	data, err := r.client.Get(gopts.Context, string(key)).Bytes()
	if err != nil {
		logger.ErrorfWithContext(gopts.Context, `[redis](%s) get failed. %s`, key, err.Error())
		output.err = err
		return output
	}

	ttl, err := r.client.TTL(gopts.Context, string(key)).Result()
	if err != nil {
		logger.ErrorfWithContext(gopts.Context, `[redis](%s) get.ttl failed. %s`, key, err.Error())
		output.err = err
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
		output.err = ErrRedisNotReady
		return output
	}

	// 调用 redis.Del() 进行删除
	// 注意：redis 的删除策略为惰性删除，并不确保立即删除，并且删除键值对会占用 CPU 资源，尤其是大量删除时
	// if err := r.client.Del(dopts.Context, string(key)).Err(); err != nil {
	// 	logger.ErrorfWithContext(dopts.Context, `[redis](%s) del failed. %s`, key, err.Error())
	// }

	// 使用 redis.RenameNX() 后设置过期时间的方式进行平滑删除
	// 相较于 redis.Del()，定时删除可以在一定程度上分摊删除操作的 CPU 负载
	_, output.err = r.client.TxPipelined(dopts.Context, func(pipe redis.Pipeliner) error {
		newkey := delSuffixKey(key)
		pipe.RenameNX(dopts.Context, string(key), string(newkey))
		// 将 key 的过期时间(删除时间)设为 0-3s, 防止集中删除
		pipe.PExpire(dopts.Context, string(newkey), time.Duration(rand.Intn(3000))*time.Millisecond)
		return nil
	})

	if output.err != nil {
		logger.ErrorfWithContext(dopts.Context, `[redis](%s) del failed. %s`, key, output.err.Error())
	}

	return output
}

// Expire .
func (r *rkv) Expire(key keyword, opts ...func(o *ExpireOptions)) *StatusOutput {
	var eopts = expireoptions(opts...)

	output := &StatusOutput{baseOutput: baseOutput{ctx: eopts.Context, key: string(key)}}
	if !r.isReady() {
		output.err = ErrRedisNotReady
		return output
	}

	if !eopts.Expiry.IsZero() {
		// ExpireAt
		output.err = r.client.PExpireAt(eopts.Context, string(key), eopts.Expiry).Err()
	} else {
		// TTL
		output.err = r.client.PExpire(eopts.Context, string(key), eopts.TTL).Err()
	}

	return output
}

// Mutex .
func (r *rkv) Mutex(key keyword, opts ...func(o *MutexOptions)) *Mutex {
	if !r.isReady() {
		return &Mutex{err: ErrRedisNotReady}
	}

	var mopts = mutexoptions(opts...)
	return &Mutex{
		rkv:  r,
		key:  lockPrefixKey(string(key)),
		opts: mopts,
		tk:   time.NewTicker(mopts.Heartbeat),
		tm:   time.NewTimer(mopts.TTL),
	}
}

// IsExists .
func (r *rkv) IsExists(key string, opts ...func(o *GetOptions)) *BoolOutput {
	var gopts = getoptions(opts...)

	output := &BoolOutput{baseOutput: baseOutput{ctx: gopts.Context, key: string(key)}}
	if !r.isReady() {
		output.err = ErrRedisNotReady
		return output
	}

	if r.client.Exists(gopts.Context, key).Val() != 0 {
		output.val = true
	}
	return output
}

func (r *rkv) Close() error {
	if r.closing != nil {
		return r.closing()
	}
	return nil
}

// ping .
func (r *rkv) ping() error {
	if err := r.client.Ping(r.opts.Context).Err(); err != nil {
		logger.ErrorfWithContext(r.opts.Context, `[redis] ping failed. %s`, err.Error())
		r.active = false
		return err
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

	client, close := options.Cmdable(options)

	r := &rkv{
		id:      id,
		opts:    options,
		client:  client,
		closing: close,
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
