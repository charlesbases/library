package logger

import (
	"os"
	"os/signal"

	"library"

	"github.com/cihub/seelog"
)

// logger .
type logger struct {
	opts   *Options
	logger seelog.LoggerInterface
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
	if l.opts.Service != "" {
		l.opts.Service = " " + "[" + l.opts.Service + "]"
	}
	if l.opts.Filename == "" {
		l.opts.Filename = DefaultFilename
	}

	logger, _ := seelog.LoggerFromConfigAsBytes([]byte(`
			<?xml version="1.0" encoding="utf-8" ?>
			<seelog levels="trace,debug,info,warn,error,critical">
				<outputs formatid="main">
					<filter levels="trace">
						<console formatid="main"/>
					</filter>
					<filter levels="debug">
						<console formatid="debug"/>
					</filter>
					<filter levels="info">
						<console formatid="info"/>
					</filter>
					<filter levels="warn">
						<console formatid="main"/>
					</filter>
					<filter levels="error,critical">
						<console formatid="error"/>
					</filter>
					<rollingfile formatid="main" type="date" filename="` + l.opts.Filename + `" datepattern="2006-01-02" maxrolls="30" namemode="prefix"/>
				</outputs>
				<formats>
					<format id="main"  format="[%Date(2006-01-02 15:04:05.000)] [%LEV]` + l.opts.Service + ` %Func %File:%Line ==&gt; %Msg%n"/>
					<format id="info"  format="%EscM(32)[%Date(2006-01-02 15:04:05.000)] [%LEV]` + l.opts.Service + ` %Func %File:%Line ==&gt; %Msg%n%EscM(0)"/>
					<format id="debug" format="%EscM(35)[%Date(2006-01-02 15:04:05.000)] [%LEV]` + l.opts.Service + ` %Func %File:%Line ==&gt; %Msg%n%EscM(0)"/>
					<format id="error" format="%EscM(31)[%Date(2006-01-02 15:04:05.000)] [%LEV]` + l.opts.Service + ` %Func %File:%Line ==&gt; %Msg%n%EscM(0)"/>
				</formats>
			</seelog>`))
	logger.SetAdditionalStackDepth(l.opts.Skip)
	l.logger = logger

	go l.flush()
}

// Trace .
func (l *logger) Trace(v ...interface{}) {
	l.logger.Trace(v...)
}

// Tracef .
func (l *logger) Tracef(format string, params ...interface{}) {
	l.logger.Tracef(format, params...)
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
	l.logger.Critical(v...)
}

// Fatalf .
func (l *logger) Fatalf(format string, params ...interface{}) {
	l.logger.Criticalf(format, params...)
}

// flush .
func (l *logger) flush() {
	s := make(chan os.Signal)
	signal.Notify(s, library.Shutdown()...)
	<-s
	l.logger.Flush()
	l.logger.Close()
}
