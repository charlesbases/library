package webmail

type Client interface {
	// Send 发送消息
	Send(v ...Message) error
}

type Message interface {
}
