package orm

import (
	"library/database"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// defaultMySQL default mysql dsn
const defaultMySQL = "root:123456@tcp(127.0.0.1:3306)/user?charset=utf8mb4&parseTime=True&loc=Local"

// MySQL .
var MySQL *mysqlDialer

type mysqlDialer struct{}

// Dialector .
func (d *mysqlDialer) Dialector(opts *database.Options) gorm.Dialector {
	if len(opts.Addrs) == 0 {
		opts.Addrs = []string{defaultMySQL}
	}
	return mysql.Open(opts.Addrs[0])
}

// Name .
func (d *mysqlDialer) Name() string {
	return "MySQL"
}
