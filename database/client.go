package database

import (
	"time"

	"github.com/pkg/errors"
)

var (
	// ErrorInvaildDsn invalid addrs
	ErrorInvaildDsn = errors.New("invalid dsn of database")
	// ErrorDatabaseNil db is not initialized or closed
	ErrorDatabaseNil = errors.New("database is not ready")
)

const (
	// DefaultMaxIdleConns default of db.MaxIdleConns
	DefaultMaxIdleConns = 200
	// DefaultConnMaxIdleTime default of db.ConnMaxIdleTime
	DefaultConnMaxIdleTime = 3 * time.Second
)

// Options .
type Options struct {
	// Address 连接地址
	Address string
	// MaxIdleConns 连接池空闲连接数。default: DefaultMaxIdleConns
	MaxIdleConns int
	// ConnMaxIdleTime 连接最大空闲时间，超时则清理。default: DefaultConnMaxIdleTime
	ConnMaxIdleTime time.Duration
	// MaxOpenConns 连接池最大连接数
	// <= 0 means unlimited
	MaxOpenConns int
}
