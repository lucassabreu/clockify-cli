package api

// Logger for the Client
type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// SetDebugLogger debug logger
func (c *client) SetDebugLogger(logger Logger) Client {
	c.debugLogger = logger
	return c
}

func (c *client) debugf(format string, v ...interface{}) {
	if c.debugLogger == nil {
		return
	}

	c.debugLogger.Printf(format, v...)
}

// SetInfoLogger info logger
func (c *client) SetInfoLogger(logger Logger) Client {
	c.infoLogger = logger
	return c
}

func (c *client) infof(format string, v ...interface{}) {
	if c.infoLogger == nil {
		return
	}

	c.infoLogger.Printf(format, v...)
}
