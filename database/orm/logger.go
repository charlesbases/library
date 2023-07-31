package orm

import (
	"context"
	"fmt"
	"time"

	"github.com/charlesbases/logger"
	glogger "gorm.io/gorm/logger"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/sonyflake"
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

// traceid 从 context.Context 中 获取 xcontext.HeaderTraceID
func (log *log) traceid(ctx context.Context) sonyflake.ID {
	if ctx == context.Background() {
		return 0
	}
	id, _ := ctx.Value(library.TraceID).(sonyflake.ID)
	return id
}

// named .
func (log *log) named(id sonyflake.ID) *logger.Logger {
	return logger.Named(id.String())
}

// warp .
func (log *log) warp(format string) string {
	return log.dt + format
}

func (log *log) LogMode(level glogger.LogLevel) glogger.Interface {
	return log
}

func (log *log) Info(ctx context.Context, format string, v ...interface{}) {
	if id := log.traceid(ctx); id != 0 {
		log.named(id).Infof(log.warp(format), v...)
	} else {
		logger.Infof(log.warp(format), v...)
	}
}

func (log *log) Warn(ctx context.Context, format string, v ...interface{}) {
	if id := log.traceid(ctx); id != 0 {
		log.named(id).Warnf(log.warp(format), v...)
	} else {
		logger.Warnf(log.warp(format), v...)
	}
}

func (log *log) Error(ctx context.Context, format string, v ...interface{}) {
	if id := log.traceid(ctx); id != 0 {
		log.named(id).Errorf(log.warp(format), v...)
	} else {
		logger.Errorf(log.warp(format), v...)
	}
}

func (log *log) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	traceid := log.traceid(ctx)

	if err != nil /* && !errors.Is(err, gorm.ErrRecordNotFound)*/ {
		if traceid != 0 {
			log.named(traceid).Errorf(log.warp("%s | %s"), sql, err.Error())
		} else {
			logger.Errorf(log.warp("%s | %s"), sql, err.Error())
		}
	} else {
		if traceid != 0 {
			log.named(traceid).Debugf(log.warp("%s | %d rows | %v"), sql, rows, elapsed)
		} else {
			logger.Debugf(log.warp("%s | %d rows | %v"), sql, rows, elapsed)
		}
	}
}
