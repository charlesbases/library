package orm

import (
	"context"
	"fmt"
	"time"

	glogger "gorm.io/gorm/logger"

	"github.com/charlesbases/logger"
)

// log .
type log struct {
	// dt driver.Driver.Type
	dt string
}

// custom .
func custom(dt string) glogger.Interface {
	return &log{dt: fmt.Sprintf("[%s] ", dt)}
}

// LogMode .
func (log *log) LogMode(level glogger.LogLevel) glogger.Interface {
	return log
}

// Info .
func (log *log) Info(ctx context.Context, format string, v ...interface{}) {
	logger.WithContext(ctx).Infof(log.dt+format, v...)
}

// Warn .
func (log *log) Warn(ctx context.Context, format string, v ...interface{}) {
	logger.WithContext(ctx).Warnf(log.dt+format, v...)
}

// Error .
func (log *log) Error(ctx context.Context, format string, v ...interface{}) {
	logger.WithContext(ctx).Errorf(log.dt+format, v...)
}

// Trace .
func (log *log) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil /* && !errors.Is(err, gorm.ErrRecordNotFound)*/ {
		logger.CallerSkip(5).WithContext(ctx).Errorf("%s%s | %v", log.dt, sql, err)
	} else {
		logger.CallerSkip(3).WithContext(ctx).Debugf("%s%s | %d rows | %v", log.dt, sql, rows, elapsed)
	}
}
