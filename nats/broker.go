package nats

type Broker interface {
	Options() *Options
	Address() string
	Connect() error
	Disconnect() error
	Publish(topic string, message *Message, opts ...PublishOption) error
	Subscribe(topic string, handler Handler, opts ...SubscribeOption) error
	String() string
}

type Header map[string]string

type Handler func(event Event) error

// Message .
type Message struct {
	Header Header
	Data   interface{}
}

type Event interface {
	Topic() string
	Reply() string

	// Body return bytes of Message
	Body() []byte
	// Header return Message.Header
	Header() Header
	// Respond allows a convenient way to respond to requests in service based subscriptions
	Respond(v interface{}) error

	// Unmarshal unmarshal for Message.Data.
	Unmarshal(pointer interface{}) error
}
