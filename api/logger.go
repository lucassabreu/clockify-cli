package api

// Logger for the Client
type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// SetDebugLogger debug logger
func (c *Client) SetDebugLogger(logger Logger) *Client {
	c.debugLogger = logger
	return c
}

func (c *Client) debugf(format string, v ...interface{}) {
	if c.debugLogger == nil {
		return
	}

	c.debugLogger.Printf(format, v...)
}
