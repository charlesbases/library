package orm

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	logger "library/logger/seelog"

	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var (
	once      sync.Once
	defaultDB *db
)

// db .
type db struct {
	gormDB *gorm.DB
}

// Init init defaultDB
func Init(fn Dialector, opts ...Option) {
	once.Do(func() {
		var options = DefaultOptions()
		for _, opt := range opts {
			opt(options)
		}

		gormDB, err := gorm.Open(fn.Dialector(options), &gorm.Config{Logger: gormlogger(fn.Name(), options.ShowSQL)})
		if err != nil {
			logger.Fatalf("database connect failed. %v", err)
		}

		defaultDB = &db{gormDB: gormDB}
	})
}

// New new db
func New(fn Dialector, opts ...Option) *db {
	var options = DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	gormDB, err := gorm.Open(fn.Dialector(options), &gorm.Config{Logger: gormlogger(fn.Name(), options.ShowSQL)})
	if err != nil {
		logger.Fatalf("database connect failed. %v", err)
	}

	return &db{gormDB: gormDB}
}

// DB .
func DB() *gorm.DB {
	return defaultDB.gormDB
}

// Transaction .
func (db *db) Transaction(fs ...func(ts *gorm.DB) error) error {
	ts := db.gormDB.Begin()
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

type Dialector interface {
	Dialector(optons *Options) gorm.Dialector
	Name() string
}

// l .
type l struct {
	showsql    bool
	namespaces string
}

// gormlogger .
func gormlogger(name string, sql bool) glogger.Interface {
	return &l{namespaces: fmt.Sprintf("[%s] >>> ", name), showsql: sql}
}

func (l *l) LogMode(level glogger.LogLevel) glogger.Interface {
	return l
}

func (l *l) Info(ctx context.Context, s string, i ...interface{}) {
	logger.Infof(l.namespaces+s, i...)
}

func (l *l) Warn(ctx context.Context, s string, i ...interface{}) {
	logger.Warnf(l.namespaces+s, i...)
}

func (l *l) Error(ctx context.Context, s string, i ...interface{}) {
	logger.Errorf(l.namespaces+s, i...)
}

func (l *l) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.showsql {
		elapsed := time.Since(begin)
		sql, rows := fc()

		switch {
		case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
			logger.Errorf(l.namespaces+"%s | %v", sql, err)
		default:
			logger.Debugf(l.namespaces+"%s | %d rows | %v", sql, rows, elapsed)
		}
	}
}
