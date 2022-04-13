package orm

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

func TestMysql(t *testing.T) {
	Init(MySQL, ShowSQL(true))
}

func TestPostgres(t *testing.T) {
	Init(PostgreSQL, ShowSQL(true))
}

func TestTransaction(t *testing.T) {
	Init(MySQL, ShowSQL(true))

	Transaction(DB(),
		func(tx *gorm.DB) error {
			tx.Table("a").Create(nil)
			return nil
		},
		func(tx *gorm.DB) error {
			tx.Table("b").Updates(nil)
			return errors.New("transaction test")
		})
}
