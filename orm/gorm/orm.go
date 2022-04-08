package gorm

import (
	"context"
	"errors"
	"sync"
	"time"

	logger "library/logger/seelog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

const dsn = "root:123456@tcp(192.168.10.188:3306)/wallet?charset=utf8mb4&parseTime=True&loc=Local"

var db *gorm.DB
var once sync.Once

// Init .
func Init() {
	once.Do(func() {
		gdb, err := gorm.Open(mysql.Open(dsn))
		if err != nil {
			logger.Fatalf(" - db connect(%s) error - %v", dsn, err)
		}

		gdb.Logger = new(l)

		db = gdb
	})
}

// DB .
func DB() *gorm.DB {
	if db != nil {
		return db
	}
	panic("db is not initialized")
}

// Transaction .
func Transaction(fs ...func(ts *gorm.DB) error) error {
	ts := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			ts.Rollback()
		}
	}()

	for _, f := range fs {
		if err := f(ts); err != nil {
			ts.Rollback()
			return err
		}
	}

	if err := ts.Commit().Error; err != nil {
		ts.Rollback()
		return err
	}
	return nil
}

const namespaces = "[MySQL] >>> "

// l .
type l struct {
	showsql bool
}

func (l *l) LogMode(level glogger.LogLevel) glogger.Interface {
	return l
}

func (l *l) Info(ctx context.Context, s string, i ...interface{}) {
	logger.Infof(namespaces+s, i...)
}

func (l *l) Warn(ctx context.Context, s string, i ...interface{}) {
	logger.Warnf(namespaces+s, i...)
}

func (l *l) Error(ctx context.Context, s string, i ...interface{}) {
	logger.Errorf(namespaces+s, i...)
}

func (l *l) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.showsql {
		elapsed := time.Since(begin)
		sql, rows := fc()

		switch {
		case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
			logger.Errorf(namespaces+"%s | %v", sql, err)
		default:
			logger.Debugf(namespaces+"%s | %d rows | %v", sql, rows, elapsed)
		}
	}
}
