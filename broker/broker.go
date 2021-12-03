package broker

type Broker interface {
	Init(opts ...Option) error
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

// NewHander .
func NewHander() Header {
	return make(map[string]string)
}

type Event interface {
	Topic() string
	Reply() string

	// Body return bytes of Message
	Body() []byte
	// Header return Message.Header
	Header() Header

	// Unmarshal unmarshal for Message.Data.
	Unmarshal(pointer interface{}) error
}
