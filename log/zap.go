package log

import (
	"os"
	"os/signal"
	"sync"

	"library"

	"go.uber.org/zap"
	_ "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "go.uber.org/zap/zapcore"
)

// logZap .
type logZap struct {
	opts *Options
	once sync.Once

	logger *zap.SugaredLogger
}

// NewZap .
func NewZap(opts ...Option) logger {
	logger := new(logZap)
	logger.configure(opts...)
	return logger
}

// configure .
func (log *logZap) configure(opts ...Option) {
	var options = new(Options)
	for _, opt := range opts {
		opt(options)
	}
	log.opts = options

	if log.opts.Skip == 0 {
		log.opts.Skip = DefaultSkip
	}
	if log.opts.Filename == "" {
		log.opts.Filename = DefaultFilename
	}
	if log.opts.DateFormat == "" {
		log.opts.DateFormat = DefaultDateFormat
	}

	// 编码器
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("[" + log.opts.DateFormat + "]")
	cfg.EncodeLevel = log.color
	cfg.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.ConsoleSeparator = " "
	encoder := zapcore.NewConsoleEncoder(cfg)

	// 控制台输出
	console := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return true
	})

	logger := zap.New(
		zapcore.NewTee(zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), console)),
		zap.AddCaller(),
		zap.AddCallerSkip(log.opts.Skip),
	)
	log.logger = logger.Named("[" + log.opts.Service + "]").Sugar()

	go log.flush()
}

// color .
func (log *logZap) color(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString(_colorString[LEVEL_DEBUG])
	case zapcore.InfoLevel:
		enc.AppendString(_colorString[LEVEL_INFO])
	case zapcore.WarnLevel:
		enc.AppendString(_colorString[LEVEL_WARN])
	case zapcore.ErrorLevel:
		enc.AppendString(_colorString[LEVEL_ERROR])
	case zapcore.DPanicLevel:
		enc.AppendString(_colorString[LEVEL_FATAL])
	case zapcore.PanicLevel:
		enc.AppendString(_colorString[LEVEL_FATAL])
	case zapcore.FatalLevel:
		enc.AppendString(_colorString[LEVEL_FATAL])
	default:
		enc.AppendString(_colorString[LEVEL_TRACE])
	}
}

// Trace .
func (log *logZap) Trace(v ...interface{}) {
	// zap does not support trace
}

// Tracef .
func (log *logZap) Tracef(format string, params ...interface{}) {
	// zap does not support trace
}

// Debug .
func (log *logZap) Debug(v ...interface{}) {
	log.logger.Debug(v...)
}

// Debugf .
func (log *logZap) Debugf(format string, params ...interface{}) {
	log.logger.Debugf(format, params...)
}

// Info .
func (log *logZap) Info(v ...interface{}) {
	log.logger.Info(v...)
}

// Infof .
func (log *logZap) Infof(format string, params ...interface{}) {
	log.logger.Infof(format, params...)
}

// Warn .
func (log *logZap) Warn(v ...interface{}) {
	log.logger.Warn(v...)
}

// Warnf .
func (log *logZap) Warnf(format string, params ...interface{}) {
	log.logger.Warnf(format, params...)
}

// Error .
func (log *logZap) Error(v ...interface{}) {
	log.logger.Error(v...)
}

// Errorf .
func (log *logZap) Errorf(format string, params ...interface{}) {
	log.logger.Errorf(format, params...)
}

// Fatal .
func (log *logZap) Fatal(v ...interface{}) {
	log.logger.Fatal(v...)
}

// Fatalf .
func (log *logZap) Fatalf(format string, params ...interface{}) {
	log.logger.Fatalf(format, params...)
}

// String .
func (log *logZap) String() string {
	return "zap"
}

// flush .
func (log *logZap) flush() {
	log.once.Do(func() {
		s := make(chan os.Signal)
		signal.Notify(s, library.Shutdown()...)
		<-s
		log.logger.Sync()
	})
}
