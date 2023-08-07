package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/charlesbases/library/logger"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/codec"
)

// lockprefix redis 分布式锁的 key 前缀
var lockprefix = KeyPrefix("lock_")

// Key 添加 key 的前缀或后缀
type Key func(key string) string

// Key .
func (k Key) Key(key string) string {
	return k(key)
}

// KeyPrefix .
func KeyPrefix(prefix string) Key {
	return func(key string) string {
		return prefix + key
	}
}

// KeySuffix .
func KeySuffix(suffix string) Key {
	return func(key string) string {
		return key + suffix
	}
}

// Mutex redis 分布式锁
type Mutex struct {
	*rkv
	key  string
	opts *MutexOptions

	tk *time.Ticker // 尝试获取锁的时间间隔
	tm *time.Timer  // 超时后自动删除
}

// Lock .
func (m *Mutex) Lock() {
	for {
		select {
		case <-m.tk.C:
			ok, _ := m.client.SetNX(m.opts.Context, m.key, m.opts.Mark, m.opts.TTL).Result()
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

	m.Del(m.key, func(o *DelOptions) {
		o.Context = m.opts.Context
	})
}

// Input .
type Input struct {
	Key string      `json:"key"`
	Val interface{} `json:"val"`
}

type Output interface {
	Key() string
	RawMessage() []byte
	TTL() time.Duration
	Expiry() time.Time
	Unmarshal(v interface{}) error
}

// output .
type output struct {
	key       string
	data      []byte
	ttl       time.Duration
	expiry    time.Time
	marshaler codec.Marshaler
}

func (o *output) Key() string {
	return o.key
}

func (o *output) RawMessage() []byte {
	return o.data
}

func (o *output) TTL() time.Duration {
	return o.ttl
}

func (o *output) Expiry() time.Time {
	return o.expiry
}

func (o *output) Unmarshal(v interface{}) error {
	if len(o.data) == 0 {
		return fmt.Errorf("[redis](%s) value is nil.", o.key)
	}

	return o.marshaler.Unmarshal(o.data, v)
}

// r default client
var r *rkv

// rkv .
type rkv struct {
	opts    *Options
	client  redis.Cmdable
	closing func() error
}

// Set .
func (r *rkv) Set(rec *Input, opts ...func(o *SetOptions)) error {
	var sopts = parsesetoptions(opts...)

	if !sopts.Expiry.IsZero() {
		sopts.TTL = time.Since(sopts.Expiry)
	}

	data, err := sopts.Marshaler.Marshal(rec.Val)
	if err != nil {
		logger.ErrorfWithContext(sopts.Context, `[redis](%s) set failed. %s`, rec.Key, err.Error())
		return err
	}

	if err := r.client.Set(sopts.Context, rec.Key, data, sopts.TTL).Err(); err != nil {
		logger.ErrorfWithContext(sopts.Context, `[redis](%s) set failed. %s`, rec.Key, err.Error())
		return err
	}

	return nil
}

// Get .
func (r *rkv) Get(key string, opts ...func(o *GetOptions)) (Output, error) {
	var gopts = parsegetoptions(opts...)

	data, err := r.client.Get(gopts.Context, key).Bytes()
	if err != nil {
		logger.ErrorfWithContext(gopts.Context, `[redis](%s) get failed. %s`, key, err.Error())
		return &output{key: key}, err
	}

	ttl, err := r.client.TTL(gopts.Context, key).Result()
	if err != nil {
		logger.ErrorfWithContext(gopts.Context, `[redis](%s) get.ttl failed. %s`, key, err.Error())
		return &output{key: key}, err
	}

	return &output{
		key:       key,
		data:      data,
		ttl:       ttl,
		expiry:    time.Now().Add(ttl),
		marshaler: gopts.Marshaler,
	}, nil
}

// Del .
func (r *rkv) Del(key string, opts ...func(o *DelOptions)) error {
	var dopts = parsedeloptions(opts...)

	newkey := fmt.Sprintf(`%s_delete_%d`, key, library.NowDuration())
	// rename
	if err := r.client.RenameNX(dopts.Context, key, newkey).Err(); err != nil {
		logger.ErrorfWithContext(dopts.Context, `[redis](%s) del failed. %s`, key, err.Error())
		return err
	}
	// delete
	go func() {
		if err := r.client.Del(dopts.Context, newkey).Err(); err != nil {
			logger.ErrorfWithContext(dopts.Context, `[redis](%s) del failed. %s`, key, err.Error())
		}
	}()

	return nil
}

// Expire .
func (r *rkv) Expire(key string, opts ...func(o *ExpireOptions)) error {
	var eopts = parseexpireoptions(opts...)

	// ExpireAt
	if !eopts.Expiry.IsZero() {
		return r.client.ExpireAt(eopts.Context, key, eopts.Expiry).Err()
	}

	// TTL
	return r.client.Expire(eopts.Context, key, eopts.TTL).Err()
}

// Mutex .
func (r *rkv) Mutex(key string, opts ...func(o *MutexOptions)) *Mutex {
	var mopts = parsemutexoptions(opts...)
	return &Mutex{
		rkv:  r,
		key:  lockprefix.Key(key),
		opts: mopts,
		tk:   time.NewTicker(mopts.Heartbeat),
		tm:   time.NewTimer(mopts.TTL),
	}
}

// IsExists .
func (r *rkv) IsExists(key string, opts ...func(o *GetOptions)) bool {
	var gopts = parsegetoptions(opts...)
	if r.client.Exists(gopts.Context, key).Val() != 0 {
		return true
	}
	return false
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
		return err
	}
	return nil
}

// disprint .
type disprint struct{}

// Printf .
func (l *disprint) Printf(_ context.Context, _ string, _ ...interface{}) {}

// NewClient .
func NewClient(opts ...func(o *Options)) (*rkv, error) {
	return parseoptions(opts...).newc()
}

// Init .
func Init(opts ...func(o *Options)) error {
	if client, err := parseoptions(opts...).newc(); err != nil {
		return err
	} else {
		r = client
	}
	return nil
}

// Client .
func Client() *rkv {
	if r != nil {
		return r
	}
	logger.Fatal(`redis is nil`)
	return nil
}

// Close .
func Close() error {
	if r != nil {
		return r.Close()
	}
	return nil
}
