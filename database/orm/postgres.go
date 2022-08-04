package orm

import (
	"library/database"

	"github.com/charlesbases/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Postgres Postgres
var Postgres *postgresDialer

type postgresDialer struct{}

// Dialector .
func (d *postgresDialer) Dialector(opts *database.Options) gorm.Dialector {
	if len(opts.Addrs) == 0 {
		logger.Fatal(database.ErrorInvaildAddrs)
	}
	return postgres.Open(opts.Addrs[0])
}

// Name .
func (d *postgresDialer) Name() string {
	return "PostgreSQL"
}
