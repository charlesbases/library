package logger

var log logger

type logger interface {
	Trace(v ...interface{})
	Tracef(format string, params ...interface{})
	Debug(v ...interface{})
	Debugf(format string, params ...interface{})
	Info(v ...interface{})
	Infof(format string, params ...interface{})
	Warn(v ...interface{})
	Warnf(format string, params ...interface{})
	Error(v ...interface{})
	Errorf(format string, params ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, params ...interface{})
	Log(v ...interface{}) error
	String() string
}

// Trace .
func Trace(v ...interface{}) {
	log.Trace(v...)
}

// Tracef .
func Tracef(format string, params ...interface{}) {
	log.Tracef(format, params...)
}

// Debug .
func Debug(v ...interface{}) {
	log.Debug(v...)
}

// Debugf .
func Debugf(format string, params ...interface{}) {
	log.Debugf(format, params...)
}

// Info .
func Info(v ...interface{}) {
	log.Info(v...)
}

// Infof .
func Infof(format string, params ...interface{}) {
	log.Infof(format, params...)
}

// Warn .
func Warn(v ...interface{}) {
	log.Warn(v...)
}

// Warnf .
func Warnf(format string, params ...interface{}) {
	log.Warnf(format, params...)
}

// Error .
func Error(v ...interface{}) {
	log.Error(v...)
}

// Errorf .
func Errorf(format string, params ...interface{}) {
	log.Errorf(format, params...)
}

// Fatal .
func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

// Fatalf .
func Fatalf(format string, params ...interface{}) {
	log.Fatalf(format, params...)
}
