package redis

import (
	"context"
	"time"

	"github.com/charlesbases/library/codec"
)

// baseOutput .
type baseOutput struct {
	ctx context.Context
	key string
	err error
}

// Ctx .
func (o *baseOutput) Ctx() context.Context {
	return o.ctx
}

// Key .
func (o *baseOutput) Key() string {
	return o.key
}

// Err .
func (o *baseOutput) Err() error {
	return o.err
}

// StatusOutput .
type StatusOutput struct {
	baseOutput
}

// BoolOutput .
type BoolOutput struct {
	baseOutput

	val bool
}

// Val .
func (o *BoolOutput) Val() bool {
	return o.val
}

// BytesOutput .
type BytesOutput struct {
	baseOutput

	val       []byte
	ttl       time.Duration
	expiry    time.Time
	marshaler codec.Marshaler
}

// Val .
func (o *BytesOutput) Val() []byte {
	return o.val
}

// TTL .
func (o *BytesOutput) TTL() time.Duration {
	return o.ttl
}

// Expiry .
func (o *BytesOutput) Expiry() time.Time {
	return o.expiry
}

// Unmarshal .
func (o *BytesOutput) Unmarshal(v interface{}) error {
	return o.marshaler.Unmarshal(o.val, v)
}
