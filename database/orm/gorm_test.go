package orm

import (
	"context"
	"testing"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/database"
	"github.com/charlesbases/library/database/orm/driver"
	"github.com/charlesbases/library/sonyflake"
)

func TestGorm(t *testing.T) {
	db, err := New(new(driver.Postgres), func(o *database.Options) {
		o.Address = "host=10.64.10.210 port=32537 user=postgres password=mxpostgres dbname=auth sslmode=disable TimeZone=Asia/Shanghai"
	})
	if err != nil {
		panic(err)
	}

	var count int64
	// With TraceID
	{
		ctx := context.WithValue(context.Background(), library.TraceID, sonyflake.NextID())
		if err := db.WithContext(ctx).Table("user_info").Count(&count).Error; err != nil {
			panic(err)
		}
	}

	{
		if err := db.Table("user_info").Count(&count).Error; err != nil {
			panic(err)
		}
	}
}
