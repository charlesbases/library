package database

import "github.com/pkg/errors"

var (
	// ErrorInvaildDsn invalid addrs
	ErrorInvaildDsn = errors.New("invalid dsn of database")
	// ErrorDatabaseNil db is not initialized or closed
	ErrorDatabaseNil = errors.New("database is not ready")
)

// DefaultMaxIdleConns default of MaxIdleConns
const DefaultMaxIdleConns = 1000

// Options .
type Options struct {
	// Address 连接地址
	Address string
	// MaxIdleConns 连接池空闲连接数
	// default 1000
	MaxIdleConns int
	// MaxOpenConns 连接池最大连接数
	// <= 0 means unlimited
	MaxOpenConns int
}
