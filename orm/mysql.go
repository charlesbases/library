package orm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// defaultMySQL default mysql dsn
const defaultMySQL = "root:123456@tcp(192.168.10.188:3306)/wallet?charset=utf8mb4&parseTime=True&loc=Local"

// MySQL .
var MySQL *mysqlDialector

type mysqlDialector struct{}

// Dialector .
func (p *mysqlDialector) Dialector(opts *Options) gorm.Dialector {
	if len(opts.Addrs) == 0 {
		opts.Addrs = []string{defaultPostgresDSN}
	}
	return mysql.Open(opts.Addrs[0])
}

// Name .
func (p *mysqlDialector) Name() string {
	return "MySQL"
}
