package orm

import (
	"context"
	"fmt"
	"time"

	glogger "gorm.io/gorm/logger"

	"github.com/charlesbases/library/logger"
)

// log .
type log struct {
	// dt driver.Driver.Type
	dt string
}

// custom .
func custom(dt string) glogger.Interface {
	return &log{dt: fmt.Sprintf("[%s] >>> ", dt)}
}

// warp .
func (log *log) warp(format string) string {
	return log.dt + format
}

func (log *log) LogMode(level glogger.LogLevel) glogger.Interface {
	return log
}

func (log *log) Info(ctx context.Context, format string, v ...interface{}) {
	logger.InfofWithContext(ctx, log.warp(format), v...)
}

func (log *log) Warn(ctx context.Context, format string, v ...interface{}) {
	logger.WarnfWithContext(ctx, log.warp(format), v...)
}

func (log *log) Error(ctx context.Context, format string, v ...interface{}) {
	logger.ErrorfWithContext(ctx, log.warp(format), v...)
}

func (log *log) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil /* && !errors.Is(err, gorm.ErrRecordNotFound)*/ {
		logger.ErrorfWithContext(ctx, log.warp("%s | %s"), sql, err.Error())
	} else {
		logger.DebugfWithContext(ctx, log.warp("%s | %d rows | %v"), sql, rows, elapsed)
	}
}
