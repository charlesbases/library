package database

import "errors"

var (
	// ErrorInvaildDsn invalid addrs
	ErrorInvaildDsn = errors.New("database: invalid dsn")
	// ErrorDatabaseNil db is not initialized or closed
	ErrorDatabaseNil = errors.New("database: db is not initialized or closed")
)

// Options .
type Options struct {
	// Address 连接地址
	Address string
	// MaxIdleConns 连接池空闲连接数
	MaxIdleConns int
	// MaxOpenConns 连接池最大连接数
	MaxOpenConns int
}
