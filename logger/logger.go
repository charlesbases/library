package logger

import (
	"context"

	"github.com/charlesbases/logger"

	"github.com/charlesbases/library"
)

// Debug .
func Debug(v ...interface{}) {
	logger.CallerSkip(1).Debug(v...)
}

// DebugWithContext .
func DebugWithContext(ctx context.Context, v ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Debug(v...)
	} else {
		logger.CallerSkip(1).Debug(v...)
	}
}

// Debugf .
func Debugf(format string, params ...interface{}) {
	logger.CallerSkip(1).Debugf(format, params...)
}

// DebugfWithContext .
func DebugfWithContext(ctx context.Context, format string, params ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Debugf(format, params...)
	} else {
		logger.CallerSkip(1).Debugf(format, params...)
	}
}

// Info .
func Info(v ...interface{}) {
	logger.CallerSkip(1).Info(v...)
}

// InfoWithContext .
func InfoWithContext(ctx context.Context, v ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Info(v...)
	} else {
		logger.CallerSkip(1).Info(v...)
	}
}

// Infof .
func Infof(format string, params ...interface{}) {
	logger.CallerSkip(1).Infof(format, params...)
}

// InfofWithContext .
func InfofWithContext(ctx context.Context, format string, params ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Infof(format, params...)
	} else {
		logger.CallerSkip(1).Infof(format, params...)
	}
}

// Warn .
func Warn(v ...interface{}) {
	logger.CallerSkip(1).Warn(v...)
}

// WarnWithContext .
func WarnWithContext(ctx context.Context, v ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Warn(v...)
	} else {
		logger.CallerSkip(1).Warn(v...)
	}
}

// Warnf .
func Warnf(format string, params ...interface{}) {
	logger.CallerSkip(1).Warnf(format, params...)
}

// WarnfWithContext .
func WarnfWithContext(ctx context.Context, format string, params ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Warnf(format, params...)
	} else {
		logger.CallerSkip(1).Warnf(format, params...)
	}
}

// Error .
func Error(v ...interface{}) {
	logger.CallerSkip(1).Error(v...)
}

// ErrorWithContext .
func ErrorWithContext(ctx context.Context, v ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Error(v...)
	} else {
		logger.CallerSkip(1).Error(v...)
	}
}

// Errorf .
func Errorf(format string, params ...interface{}) {
	logger.CallerSkip(1).Errorf(format, params...)
}

// ErrorfWithContext .
func ErrorfWithContext(ctx context.Context, format string, params ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Errorf(format, params...)
	} else {
		logger.CallerSkip(1).Errorf(format, params...)
	}
}

// Fatal .
func Fatal(v ...interface{}) {
	logger.CallerSkip(1).Fatal(v...)
}

// FatalWithContext .
func FatalWithContext(ctx context.Context, v ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Fatal(v...)
	} else {
		logger.CallerSkip(1).Fatal(v...)
	}
}

// Fatalf .
func Fatalf(format string, params ...interface{}) {
	logger.CallerSkip(1).Fatalf(format, params...)
}

// FatalfWithContext .
func FatalfWithContext(ctx context.Context, format string, params ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Fatalf(format, params...)
	} else {
		logger.CallerSkip(1).Fatalf(format, params...)
	}
}
