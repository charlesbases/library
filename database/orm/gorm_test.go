package orm

import (
	"testing"

	"github.com/charlesbases/logger"
	"gorm.io/gorm"

	"github.com/charlesbases/library/database"
	"github.com/charlesbases/library/database/orm/driver"
)

func TestGorm(t *testing.T) {
	err := Init(new(driver.Postgres), func(o *database.Options) {
		o.Address = "host=10.64.10.210 port=32537 user=postgres password=mxpostgres dbname=auth application_name=test sslmode=disable TimeZone=Asia/Shanghai"
	})
	if err != nil {
		panic(err)
	}

	var count int64
	// With TraceID
	{
		Do(func(db *gorm.DB) error {
			return db.Table("user_info").Count(&count).Error
		})
	}

	Do(func(db *gorm.DB) error {
		return db.Table("user_").Count(&count).Error
	})

	Transaction(
		func(tx *gorm.DB) error {
			return tx.Table("user_info").Count(&count).Error
		},
		func(tx *gorm.DB) error {
			return tx.Table("user_info").Count(&count).Error
		},
		func(tx *gorm.DB) error {
			return tx.Table("user_info").Count(&count).Error
		},
	)
	logger.Flush()
}
