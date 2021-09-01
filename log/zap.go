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

	sugared *zap.SugaredLogger

	once sync.Once
}

// UseZap .
func UseZap(opts ...Option) {
	logger := new(logZap)
	logger.configure(opts...)

	log = logger
}

// configure .
func (log *logZap) configure(opts ...Option) {
	var options = new(Options)
	for _, opt := range opts {
		opt(options)
	}
	log.opts = options

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
		zap.AddCallerSkip(2),
	)
	log.sugared = logger.Named(log.opts.Service).Sugar()

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
	log.sugared.Debug(v...)
}

// Debugf .
func (log *logZap) Debugf(format string, params ...interface{}) {
	log.sugared.Debugf(format, params...)
}

// Info .
func (log *logZap) Info(v ...interface{}) {
	log.sugared.Info(v...)
}

// Infof .
func (log *logZap) Infof(format string, params ...interface{}) {
	log.sugared.Infof(format, params...)
}

// Warn .
func (log *logZap) Warn(v ...interface{}) {
	log.sugared.Warn(v...)
}

// Warnf .
func (log *logZap) Warnf(format string, params ...interface{}) {
	log.sugared.Warnf(format, params...)
}

// Error .
func (log *logZap) Error(v ...interface{}) {
	log.sugared.Error(v...)
}

// Errorf .
func (log *logZap) Errorf(format string, params ...interface{}) {
	log.sugared.Errorf(format, params...)
}

// Fatal .
func (log *logZap) Fatal(v ...interface{}) {
	log.sugared.Fatal(v...)
}

// Fatalf .
func (log *logZap) Fatalf(format string, params ...interface{}) {
	log.sugared.Fatalf(format, params...)
}

// String .
func (log *logZap) String() string {
	return "log.sugared"
}

// flush .
func (log *logZap) flush() {
	log.once.Do(func() {
		s := make(chan os.Signal)
		signal.Notify(s, library.Shutdown()...)
		<-s
		log.sugared.Sync()
	})
}
