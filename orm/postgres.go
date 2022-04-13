package orm

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// defaultPostgresDSN default postgresql dsn
const defaultPostgresDSN = "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=postgres sslmode=disable"

// PostgreSQL .
var PostgreSQL *postgresDialector

type postgresDialector struct{}

// Dialector .
func (p *postgresDialector) Dialector(opts *Options) gorm.Dialector {
	if len(opts.Addrs) == 0 {
		opts.Addrs = []string{defaultPostgresDSN}
	}
	return postgres.Open(opts.Addrs[0])
}

// Name .
func (p *postgresDialector) Name() string {
	return "PostgreSQL"
}
