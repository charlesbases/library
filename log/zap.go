package log

// logZap .
type logZap struct {
	opts *Options
}

// UseZap .
func UseZap(opts ...Option) {
	var options = new(Options)
	for _, opt := range opts {
		opt(options)
	}

	log := new(logZap)
	log.opts = options
	log.configure()
}

// configure .
func (log *logZap) configure() {
	if log.opts.Filename == "" {
		log.opts.Filename = DefaultFilename
	}
	if log.opts.DateFormat == "" {
		log.opts.DateFormat = DefaultDateFormat
	}
}

func (log *logZap) Trace(v ...interface{}) {
	panic("implement me")
}

func (log *logZap) Tracef(format string, params ...interface{}) {
	panic("implement me")
}

func (log *logZap) Debug(v ...interface{}) {
	panic("implement me")
}

func (log *logZap) Debugf(format string, params ...interface{}) {
	panic("implement me")
}

func (log *logZap) Info(v ...interface{}) {
	panic("implement me")
}

func (log *logZap) Infof(format string, params ...interface{}) {
	panic("implement me")
}

func (log *logZap) Warn(v ...interface{}) {
	panic("implement me")
}

func (log *logZap) Warnf(format string, params ...interface{}) {
	panic("implement me")
}

func (log *logZap) Error(v ...interface{}) {
	panic("implement me")
}

func (log *logZap) Errorf(format string, params ...interface{}) {
	panic("implement me")
}

func (log *logZap) Fatal(v ...interface{}) {
	panic("implement me")
}

func (log *logZap) Fatalf(format string, params ...interface{}) {
	panic("implement me")
}

func (log *logZap) String() string {
	panic("implement me")
}
