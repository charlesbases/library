package zap

import (
	"io"
	"os"
	"os/signal"
	"time"

	"library"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	_ "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "go.uber.org/zap/zapcore"
)

// logger .
type logger struct {
	opts   *Options
	logger *zap.SugaredLogger
}

// New .
func New(opts ...Option) *logger {
	l := new(logger)
	l.configure(opts...)
	return l
}

// configure .
func (l *logger) configure(opts ...Option) {
	var options = new(Options)
	for _, opt := range opts {
		opt(options)
	}
	l.opts = options

	if l.opts.Skip == 0 {
		l.opts.Skip = DefaultSkip
	}
	if l.opts.Filename == "" {
		l.opts.Filename = DefaultFilename
	}
	if l.opts.DateFormat == "" {
		l.opts.DateFormat = DefaultDateFormat
	}

	// 编码器
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("[" + l.opts.DateFormat + "]")
	cfg.EncodeLevel = l.color
	cfg.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.ConsoleSeparator = " "
	encoder := zapcore.NewConsoleEncoder(cfg)

	// 日志级别
	level := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return true
	})

	logger := zap.New(
		zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level), // console
			// zapcore.NewCore(encoder, zapcore.AddSync(l.writer()), level), // file
		),
		zap.AddCaller(),
		zap.AddCallerSkip(l.opts.Skip),
	)

	if len(l.opts.Service) != 0 {
		logger = logger.Named("[" + l.opts.Service + "]")
	}
	l.logger = logger.Sugar()

	go l.flush()
}

// color .
func (l *logger) color(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var level level
	switch lv {
	case zapcore.DebugLevel:
		level = LEVEL_DBG
	case zapcore.InfoLevel:
		level = LEVEL_INF
	case zapcore.WarnLevel:
		level = LEVEL_WRN
	case zapcore.ErrorLevel:
		level = LEVEL_ERR
	case zapcore.DPanicLevel:
		level = LEVEL_FAT
	case zapcore.PanicLevel:
		level = LEVEL_FAT
	case zapcore.FatalLevel:
		level = LEVEL_FAT
	default:
		level = LEVEL_TRC
	}

	enc.AppendString(level.Sprint(level.short()))
}

// Trace .
func (l *logger) Trace(v ...interface{}) {
	// zap does not support trace
}

// Tracef .
func (l *logger) Tracef(format string, params ...interface{}) {
	// zap does not support trace
}

// Debug .
func (l *logger) Debug(v ...interface{}) {
	l.logger.Debug(v...)
}

// Debugf .
func (l *logger) Debugf(format string, params ...interface{}) {
	l.logger.Debugf(format, params...)
}

// Info .
func (l *logger) Info(v ...interface{}) {
	l.logger.Info(v...)
}

// Infof .
func (l *logger) Infof(format string, params ...interface{}) {
	l.logger.Infof(format, params...)
}

// Warn .
func (l *logger) Warn(v ...interface{}) {
	l.logger.Warn(v...)
}

// Warnf .
func (l *logger) Warnf(format string, params ...interface{}) {
	l.logger.Warnf(format, params...)
}

// Error .
func (l *logger) Error(v ...interface{}) {
	l.logger.Error(v...)
}

// Errorf .
func (l *logger) Errorf(format string, params ...interface{}) {
	l.logger.Errorf(format, params...)
}

// Fatal .
func (l *logger) Fatal(v ...interface{}) {
	l.logger.Fatal(v...)
}

// Fatalf .
func (l *logger) Fatalf(format string, params ...interface{}) {
	l.logger.Fatalf(format, params...)
}

// writer .
func (l *logger) writer() io.Writer {
	clerk, err := rotatelogs.New(
		l.opts.Filename,
		rotatelogs.WithLinkName(l.opts.Filename),
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithMaxAge(time.Hour*24*7),
	)
	if err != nil {
		panic(err)
	}
	return clerk
}

// flush .
func (l *logger) flush() {
	s := make(chan os.Signal)
	signal.Notify(s, library.Shutdown()...)
	<-s
	l.logger.Sync()
}
