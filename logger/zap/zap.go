package zap

import (
	"os"
	"os/signal"

	"library"

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
	var options = defaultOption()
	for _, opt := range opts {
		opt(options)
	}
	l.opts = options

	// 编码器
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("[" + DefaultDateFormat + "]")
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
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),                                                          // console
			zapcore.NewCore(encoder, zapcore.AddSync(NewFileWriter(l.opts.FilePath, FileWriterWithTTL(l.opts.MaxRolls))), level), // file-writer
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
		level = level_debug
	case zapcore.InfoLevel:
		level = level_info
	case zapcore.WarnLevel:
		level = level_warn
	case zapcore.ErrorLevel:
		level = level_error
	case zapcore.DPanicLevel:
		level = level_fatal
	case zapcore.PanicLevel:
		level = level_fatal
	case zapcore.FatalLevel:
		level = level_fatal
	default:
		level = level_trace
	}

	enc.AppendString(level.sprint(level.short()))
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

// flush .
func (l *logger) flush() {
	s := make(chan os.Signal)
	signal.Notify(s, library.Shutdown()...)
	<-s
	l.logger.Sync()
}
