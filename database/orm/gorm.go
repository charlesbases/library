package orm

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/charlesbases/library/database"
	"github.com/charlesbases/library/database/orm/driver"
)

var db *gorm.DB

// configuration .
func configuration(opts ...func(o *database.Options)) *database.Options {
	var options = &database.Options{
		MaxIdleConns:    database.DefaultMaxIdleConns,
		ConnMaxIdleTime: database.DefaultConnMaxIdleTime,
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// NewClient new db
func NewClient(fn driver.Driver, opts ...func(o *database.Options)) (*gorm.DB, error) {
	var options = configuration(opts...)
	if len(options.Address) == 0 {
		return nil, database.ErrorInvaildDsn
	}

	gormDB, err := gorm.Open(fn.Dialer(options), &gorm.Config{Logger: custom(fn.Type())})
	if err != nil {
		return nil, errors.Wrap(err, "open")
	}

	{
		db, err := gormDB.DB()
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}

		if err := db.Ping(); err != nil {
			return nil, errors.Wrap(err, "ping")
		}

		db.SetMaxIdleConns(options.MaxIdleConns)
		db.SetConnMaxIdleTime(options.ConnMaxIdleTime)
	}

	return gormDB, nil
}

// Init init default client
func Init(fn driver.Driver, opts ...func(o *database.Options)) error {
	client, err := NewClient(fn, opts...)
	db = client
	return err
}

// Do do something with db
func Do(fn func(db *gorm.DB) error) error {
	if db == nil {
		return database.ErrorDatabaseNil
	}
	return fn(db)
}

// DB .
func DB() *gorm.DB {
	return db
}

// Transaction transaction with db
func Transaction(fs ...func(tx *gorm.DB) error) (err error) {
	if db == nil {
		return database.ErrorDatabaseNil
	}

	tx := db.Begin()
	defer func() {
		if err != nil || recover() != nil {
			tx.Rollback()
		}
	}()

	for _, f := range fs {
		if err := f(tx); err != nil {
			return err
		}
	}

	return tx.Commit().Error
}
