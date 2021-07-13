package store

import (
	"context"
	"errors"
	"time"

	"library/metadata"
)

var (
	DefaultContext = context.Background()
	ErrNotFound    = errors.New("not found")
)

type Store interface {
	Init(...Option) error
	Options() *Options
	Read(key string, opts ...ReadOption) ([]*Record, error)
	Write(r *Record, opts ...WriteOption) error
	Delete(key string, opts ...DeleteOption) error
	List(opts ...ListOption) ([]string, error)
	Close() error
	String() string
}

type Record struct {
	Key      string            `json:"key"`
	Value    []byte            `json:"value"`
	Metadata metadata.Metadata `json:"metadata"`
	TTL      time.Duration     `json:"ttl,omitempty"`
	Expiry   time.Time         `json:"expiry,omitempty"`
}
