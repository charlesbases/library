package orm

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"library/database"

	"github.com/charlesbases/logger"
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

// New new db
func New(fn Dialector, opts ...database.Option) *db {
	var options = new(database.Options)
	for _, opt := range opts {
		opt(options)
	}
	if options.Timeout == 0 {
		options.Timeout = database.DefaultTimeout
	}

	gormDB, err := gorm.Open(fn.Dialector(options), &gorm.Config{Logger: gormlogger(fn.Name(), options.ShowSQL)})
	if err != nil {
		logger.Fatalf("database connect failed. %v", err)
	}

	return &db{gormDB: gormDB}
}

// Init init defaultDB
func Init(fn Dialector, opts ...database.Option) {
	once.Do(func() {
		var options = database.DefaultOptions()
		for _, opt := range opts {
			opt(options)
		}
		if options.Timeout == 0 {
			options.Timeout = database.DefaultTimeout
		}

		gormDB, err := gorm.Open(fn.Dialector(options), &gorm.Config{Logger: gormlogger(fn.Name(), options.ShowSQL)})
		if err != nil {
			logger.Fatalf("database connect failed. %v", err)
		}

		defaultDB = &db{gormDB: gormDB}
	})
}

// DB .
func DB() *gorm.DB {
	return defaultDB.gormDB
}

// Transaction .
func Transaction(gormDB *gorm.DB, fs ...func(tx *gorm.DB) error) error {
	tx := gormDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, f := range fs {
		if err := f(tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

type Dialector interface {
	Dialector(optons *database.Options) gorm.Dialector
	Name() string
}

// l .
type l struct {
	showsql    bool
	namespaces string
}

// gormlogger .
func gormlogger(name string, show bool) glogger.Interface {
	return &l{namespaces: fmt.Sprintf("[%s] >>> ", name), showsql: show}
}

func (l *l) LogMode(level glogger.LogLevel) glogger.Interface {
	return l
}

func (l *l) Info(ctx context.Context, format string, v ...interface{}) {
	logger.Infof(l.namespaces+format, v...)
}

func (l *l) Warn(ctx context.Context, format string, v ...interface{}) {
	logger.Warnf(l.namespaces+format, v...)
}

func (l *l) Error(ctx context.Context, format string, v ...interface{}) {
	logger.Errorf(l.namespaces+format, v...)
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
