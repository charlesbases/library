package redis

import (
	"fmt"
	"time"

	"github.com/charlesbases/logger"
	"github.com/go-redis/redis/v8"

	"library/store"
)

var (
	DefaultAddress = "redis://127.0.0.1:6379"
)

const (
	kindPrefix  = "Prefix"
	kindSuffix  = "Suffix"
	kindDefault = "Default"
)

type keyKind string

// NewStore returns a redis store
func NewStore(opts ...store.Option) store.Store {
	var options *store.Options
	for _, o := range opts {
		o(options)
	}

	s := new(rkv)
	s.options = options

	if err := s.configure(); err != nil {
		logger.Fatal("redis connect error: ", err)
	}

	return s
}

type rkv struct {
	options *store.Options
	client  *redis.Client
}

func (r *rkv) configure() error {
	var redisOptions *redis.Options
	addrs := r.options.Addresses

	if len(addrs) == 0 {
		addrs = []string{DefaultAddress}
	}

	redisOptions, err := redis.ParseURL(addrs[0])
	if err != nil {
		// Backwards compatibility
		redisOptions = &redis.Options{
			Addr:     addrs[0],
			Password: "", // no password set
			DB:       0,  // use default DB
		}
	}

	if r.options.Context == nil {
		r.options.Context = store.DefaultContext
	}

	if r.options.Auth {
		redisOptions.Password = r.options.Password
	}

	r.client = redis.NewClient(redisOptions)
	return nil
}

func (r *rkv) Init(opts ...store.Option) error {
	for _, o := range opts {
		o(r.options)
	}

	return r.configure()
}

func (r *rkv) Read(key string, opts ...store.ReadOption) ([]*store.Record, error) {
	readOptions := new(store.ReadOptions)

	for _, o := range opts {
		o(readOptions)
	}

	var keys = make([]string, 1)
	if readOptions.Prefix || readOptions.Suffix {
		// Prefix
		{
			if readOptions.Prefix {
				pkeys, err := r.keys(key, kindPrefix)
				if err != nil {
					return nil, err
				}
				keys = append(keys, pkeys...)
			}
		}
		// Suffix
		{
			if readOptions.Suffix {
				skeys, err := r.keys(key, kindSuffix)
				if err != nil {
					return nil, err
				}
				keys = append(keys, skeys...)
			}
		}
	} else {
		keys[0] = key
	}

	records := make([]*store.Record, 0, len(keys))

	for _, k := range keys {
		val, err := r.client.Get(r.options.Context, k).Bytes()

		if err != nil && err == redis.Nil {
			return nil, store.ErrNotFound
		} else if err != nil {
			return nil, err
		}
		if val == nil {
			return nil, store.ErrNotFound
		}
		d, err := r.client.TTL(r.options.Context, k).Result()
		if err != nil {
			return nil, err
		}
		records = append(records, &store.Record{
			Key:    key,
			Value:  val,
			TTL:    d,
			Expiry: time.Now().Add(d),
		})
	}
	return records, nil
}

func (r *rkv) Delete(key string, opts ...store.DeleteOption) error {
	deleteOptions := new(store.DeleteOptions)
	for _, o := range opts {
		o(deleteOptions)
	}

	var keys = make([]string, 1)
	if deleteOptions.Prefix || deleteOptions.Suffix {
		// Prefix
		{
			if deleteOptions.Prefix {
				pkeys, err := r.keys(key, kindPrefix)
				if err != nil {
					return err
				}
				keys = append(keys, pkeys...)
			}
		}
		// Suffix
		{
			if deleteOptions.Suffix {
				skeys, err := r.keys(key, kindSuffix)
				if err != nil {
					return err
				}
				keys = append(keys, skeys...)
			}
		}
	} else {
		keys[0] = key
	}

	return r.client.Del(r.options.Context, keys...).Err()
}

func (r *rkv) Write(record *store.Record, opts ...store.WriteOption) error {
	writeOptions := new(store.WriteOptions)

	for _, o := range opts {
		o(writeOptions)
	}

	if len(opts) > 0 {
		if !writeOptions.Expiry.IsZero() {
			record.TTL = time.Since(writeOptions.Expiry)
			record.Expiry = writeOptions.Expiry
		}
		if writeOptions.TTL != 0 {
			record.TTL = writeOptions.TTL
			record.Expiry = time.Now().Add(writeOptions.TTL)
		}
	}

	return r.client.Set(r.options.Context, record.Key, record.Value, record.TTL).Err()
}

func (r *rkv) List(opts ...store.ListOption) ([]string, error) {
	listOptions := new(store.ListOptions)

	for _, o := range opts {
		o(listOptions)
	}

	var pattern string
	if listOptions.Prefix != "" {
		pattern = listOptions.Prefix + "*"
	} else if listOptions.Suffix != "" {
		pattern = "*" + listOptions.Suffix
	} else {
		pattern = "*"
	}

	keys, err := r.client.Keys(r.options.Context, pattern).Result()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *rkv) Close() error {
	return r.client.Close()
}

func (r *rkv) String() string {
	return "redis"
}

func (r *rkv) Options() *store.Options {
	return r.options
}

// keys .
func (r *rkv) keys(key string, kind keyKind) ([]string, error) {
	switch kind {
	case kindPrefix:
		return r.client.Keys(r.options.Context, fmt.Sprintf(`%s*`, key)).Result()
	case kindSuffix:
		return r.client.Keys(r.options.Context, fmt.Sprintf(`*%s`, key)).Result()
	default:
		return []string{key}, nil
	}
}
