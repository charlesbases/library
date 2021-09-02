package log

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
	String() string
}
