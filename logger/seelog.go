package logger

import (
	"os"
	"os/signal"

	"library"

	"github.com/cihub/seelog"
)

// loggerSeelog .
type loggerSeelog struct {
}

// NewSeelog .
func NewSeelog() {
	see := new(loggerSeelog)
	see.configure()

	// init logger
	log = see
}

// configure .
func (log *loggerSeelog) configure() error {
	logger, _ := seelog.LoggerFromConfigAsString(`
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
					<rollingfile formatid="main" type="date" filename="./log/log.log" datepattern="2006-01-02" maxrolls="30" namemode="prefix"/>
				</outputs>
				<formats>
					<format id="main"  format="[%Date(2006-01-02 15:04:05.000)][%LEV] %Func %File:%Line ==&gt; %Msg%n"/>
					<format id="info"  format="%EscM(32)[%Date(2006-01-02 15:04:05.000)][%LEV] %Func %File:%Line ==&gt; %Msg%n%EscM(0)"/>
					<format id="debug" format="%EscM(36)[%Date(2006-01-02 15:04:05.000)][%LEV] %Func %File:%Line ==&gt; %Msg%n%EscM(0)"/>
					<format id="error" format="%EscM(31)[%Date(2006-01-02 15:04:05.000)][%LEV] %Func %File:%Line ==&gt; %Msg%n%EscM(0)"/>
				</formats>
			</seelog>`)
	logger.SetAdditionalStackDepth(1)
	seelog.UseLogger(logger)

	log.flush()
	return nil
}

// Trace .
func (log *loggerSeelog) Trace(v ...interface{}) {
	seelog.Trace(v...)
}

// Tracef .
func (log *loggerSeelog) Tracef(format string, params ...interface{}) {
	seelog.Tracef(format, params...)
}

// Debug .
func (log *loggerSeelog) Debug(v ...interface{}) {
	seelog.Debug(v...)
}

// Debugf .
func (log *loggerSeelog) Debugf(format string, params ...interface{}) {
	seelog.Debugf(format, params...)
}

// Info .
func (log *loggerSeelog) Info(v ...interface{}) {
	seelog.Info(v...)
}

// Infof .
func (log *loggerSeelog) Infof(format string, params ...interface{}) {
	seelog.Infof(format, params...)
}

// Warn .
func (log *loggerSeelog) Warn(v ...interface{}) {
	seelog.Warn(v...)
}

// Warnf .
func (log *loggerSeelog) Warnf(format string, params ...interface{}) {
	seelog.Warnf(format, params...)
}

// Error .
func (log *loggerSeelog) Error(v ...interface{}) {
	seelog.Error(v...)
}

// Errorf .
func (log *loggerSeelog) Errorf(format string, params ...interface{}) {
	seelog.Errorf(format, params...)
}

// Fatal .
func (log *loggerSeelog) Fatal(v ...interface{}) {
	seelog.Critical(v...)
}

// Fatalf .
func (log *loggerSeelog) Fatalf(format string, params ...interface{}) {
	seelog.Criticalf(format, params...)
}

// Log .
func (log *loggerSeelog) Log(v ...interface{}) error {
	seelog.Info(v...)
	return nil
}

// String .
func (log *loggerSeelog) String() string {
	return "seelog"
}

// flush .
func (log *loggerSeelog) flush() {
	go func() {
		s := make(chan os.Signal)
		signal.Notify(s, library.Shutdown()...)
		<-s
		seelog.Flush()
	}()
}
