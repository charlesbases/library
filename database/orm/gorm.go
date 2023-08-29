package orm

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/charlesbases/library/database"
	"github.com/charlesbases/library/database/orm/driver"
	"github.com/charlesbases/library/logger"
)

var db *gorm.DB

// configuration .
func configuration(opts ...func(o *database.Options)) *database.Options {
	var options = new(database.Options)
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
		return nil, fmt.Errorf("database connect failed. %s", err.Error())
	}

	{
		db, err := gormDB.DB()
		if err != nil {
			return nil, fmt.Errorf("database connect failed. %s", err.Error())
		}

		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("database ping failed. %s", err.Error())
		}

		db.SetMaxOpenConns(options.MaxOpenConns)
		db.SetMaxIdleConns(options.MaxIdleConns)
	}

	return gormDB, nil
}

// Init init default client
func Init(fn driver.Driver, opts ...func(o *database.Options)) error {
	client, err := NewClient(fn, opts...)
	if err != nil {
		return err
	}

	db = client
	return nil
}

// Do do something with db
func Do(fn func(db *gorm.DB) error) error {
	if db == nil {
		logger.Error(database.ErrorDatabaseNil.Error())
		return database.ErrorDatabaseNil
	}

	return fn(db)
}

// Transaction transaction with db
func Transaction(fs ...func(tx *gorm.DB) error) error {
	if db == nil {
		logger.Error(database.ErrorDatabaseNil.Error())
		return database.ErrorDatabaseNil
	}

	tx := db.Begin()
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
