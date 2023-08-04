package logger

import (
	"context"

	"github.com/charlesbases/logger"

	"github.com/charlesbases/library"
)

// DebugWithContext .
func DebugWithContext(ctx context.Context, v ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Debug(v...)
	} else {
		logger.CallerSkip(1).Debug(v...)
	}
}

// DebugfWithContext .
func DebugfWithContext(ctx context.Context, format string, params ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Debugf(format, params...)
	} else {
		logger.CallerSkip(1).Debugf(format, params...)
	}
}

// InfoWithContext .
func InfoWithContext(ctx context.Context, v ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Info(v...)
	} else {
		logger.CallerSkip(1).Info(v...)
	}
}

// InfofWithContext .
func InfofWithContext(ctx context.Context, format string, params ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Infof(format, params...)
	} else {
		logger.CallerSkip(1).Infof(format, params...)
	}
}

// WarnWithContext .
func WarnWithContext(ctx context.Context, v ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Warn(v...)
	} else {
		logger.CallerSkip(1).Warn(v...)
	}
}

// WarnfWithContext .
func WarnfWithContext(ctx context.Context, format string, params ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Warnf(format, params...)
	} else {
		logger.CallerSkip(1).Warnf(format, params...)
	}
}

// ErrorWithContext .
func ErrorWithContext(ctx context.Context, v ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Error(v...)
	} else {
		logger.CallerSkip(1).Error(v...)
	}
}

// ErrorfWithContext .
func ErrorfWithContext(ctx context.Context, format string, params ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Errorf(format, params...)
	} else {
		logger.CallerSkip(1).Errorf(format, params...)
	}
}

// FatalWithContext .
func FatalWithContext(ctx context.Context, v ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Fatal(v...)
	} else {
		logger.CallerSkip(1).Fatal(v...)
	}
}

// FatalfWithContext .
func FatalfWithContext(ctx context.Context, format string, params ...interface{}) {
	if id := library.ContextValueTraceID(ctx); id != 0 {
		logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 }).Fatalf(format, params...)
	} else {
		logger.CallerSkip(1).Fatalf(format, params...)
	}
}
