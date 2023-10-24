package webmail

// Client .
type Client interface {
	// Send 发送消息
	Send(v ...interface{}) error
}
