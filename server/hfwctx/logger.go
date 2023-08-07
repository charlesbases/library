package hfwctx

// Debug .
func (c *Context) Debug(v ...interface{}) {
	c.log.Debug(v...)
}

// Debugf .
func (c *Context) Debugf(format string, v ...interface{}) {
	c.log.Debugf(format, v...)
}

// Info .
func (c *Context) Info(v ...interface{}) {
	c.log.Info(v...)
}

// Infof .
func (c *Context) Infof(format string, v ...interface{}) {
	c.log.Infof(format, v...)
}

// Warn .
func (c *Context) Warn(v ...interface{}) {
	c.log.Warn(v...)
}

// Warnf .
func (c *Context) Warnf(format string, v ...interface{}) {
	c.log.Warnf(format, v...)
}

// Error .
func (c *Context) Error(v ...interface{}) {
	c.log.Error(v...)
}

// Errorf .
func (c *Context) Errorf(format string, v ...interface{}) {
	c.log.Errorf(format, v...)
}

// Fatal .
func (c *Context) Fatal(v ...interface{}) {
	c.log.Fatal(v...)
}

// Fatalf .
func (c *Context) Fatalf(format string, v ...interface{}) {
	c.log.Fatalf(format, v...)
}
