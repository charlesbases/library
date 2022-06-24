package orm

import (
	"library/database"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// defaultPostgresDSN default postgresql dsn
const defaultPostgresDSN = "host=127.0.0.1 port=5432 user=postgres password=123456 dbname=postgres sslmode=disable"

// Postgres Postgres
var Postgres *postgresDialer

type postgresDialer struct{}

// Dialector .
func (d *postgresDialer) Dialector(opts *database.Options) gorm.Dialector {
	if len(opts.Addrs) == 0 {
		opts.Addrs = []string{defaultPostgresDSN}
	}
	return postgres.Open(opts.Addrs[0])
}

// Name .
func (d *postgresDialer) Name() string {
	return "PostgreSQL"
}
