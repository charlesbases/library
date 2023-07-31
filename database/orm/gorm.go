package orm

import (
	"fmt"

	"github.com/charlesbases/logger"
	"gorm.io/gorm"

	"github.com/charlesbases/library/database"
	"github.com/charlesbases/library/database/orm/driver"
)

var db *client

// client .
type client struct {
	*gorm.DB
}

// configuration .
func configuration(opts ...func(o *database.Options)) *database.Options {
	var options = new(database.Options)
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// New new db
func New(fn driver.Driver, opts ...func(o *database.Options)) (*client, error) {
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

	return &client{DB: gormDB}, nil
}

// Init init defaultDB
func Init(fn driver.Driver, opts ...func(o *database.Options)) error {
	client, err := New(fn, opts...)
	if err != nil {
		return err
	}

	db = client
	return nil
}

// DB return db
func DB() *client {
	if db != nil {
		return db
	} else {
		logger.Fatalf("get db failed. %s", database.ErrorDatabaseNil.Error())
		return nil
	}
}

// Transaction transaction of defaultDB
func (db *client) Transaction(fs ...func(tx *gorm.DB) error) error {
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
