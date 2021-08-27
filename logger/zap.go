package logger

// loggerZap .
type loggerZap struct {
	opts *Options
}

// UseZap .
func UseZap(opts ...Option) {
	var options = new(Options)
	for _, opt := range opts {
		opt(options)
	}

	log := new(loggerZap)
	log.opts = options
	log.configure()
}

// configure .
func (log *loggerZap) configure() {
	if log.opts.Filename == "" {
		log.opts.Filename = DefaultFilename
	}
	if log.opts.DateFormat == "" {
		log.opts.DateFormat = DefaultDateFormat
	}
}

func (log *loggerZap) Trace(v ...interface{}) {
	panic("implement me")
}

func (log *loggerZap) Tracef(format string, params ...interface{}) {
	panic("implement me")
}

func (log *loggerZap) Debug(v ...interface{}) {
	panic("implement me")
}

func (log *loggerZap) Debugf(format string, params ...interface{}) {
	panic("implement me")
}

func (log *loggerZap) Info(v ...interface{}) {
	panic("implement me")
}

func (log *loggerZap) Infof(format string, params ...interface{}) {
	panic("implement me")
}

func (log *loggerZap) Warn(v ...interface{}) {
	panic("implement me")
}

func (log *loggerZap) Warnf(format string, params ...interface{}) {
	panic("implement me")
}

func (log *loggerZap) Error(v ...interface{}) {
	panic("implement me")
}

func (log *loggerZap) Errorf(format string, params ...interface{}) {
	panic("implement me")
}

func (log *loggerZap) Fatal(v ...interface{}) {
	panic("implement me")
}

func (log *loggerZap) Fatalf(format string, params ...interface{}) {
	panic("implement me")
}

func (log *loggerZap) String() string {
	panic("implement me")
}
