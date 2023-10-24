package webmail

// Client .
type Client interface {
	// Send 发送消息
	Send(v ...Message) error
}

// Message .
type Message interface {
}
