package orm

import (
	"library/database"

	"github.com/charlesbases/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySQL .
var MySQL *mysqlDialer

type mysqlDialer struct{}

// Dialector .
func (d *mysqlDialer) Dialector(opts *database.Options) gorm.Dialector {
	if len(opts.Addrs) == 0 {
		logger.Fatal(database.ErrorInvaildAddrs)
	}
	return mysql.Open(opts.Addrs[0])
}

// Name .
func (d *mysqlDialer) Name() string {
	return "MySQL"
}
