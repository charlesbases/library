package driver

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/charlesbases/library/database"
)

type Driver interface {
	Dialer(optons *database.Options) gorm.Dialector
	Type() string
}

// Mysql .
type Mysql struct{}

// Dialer .
func (*Mysql) Dialer(options *database.Options) gorm.Dialector {
	return mysql.Open(options.Address)
}

// Type .
func (*Mysql) Type() string {
	return "mysql"
}

// Postgres .
type Postgres struct{}

// Dialer .
func (*Postgres) Dialer(options *database.Options) gorm.Dialector {
	return postgres.Open(options.Address)
}

// Type .
func (*Postgres) Type() string {
	return "postgres"
}
