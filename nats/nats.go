package nats

import (
	"sync"

	"library/codec"

	"github.com/nats-io/nats.go"
)

// conn .
type conn struct {
	addrs []string
	opts  *Options

	natsConn      *nats.Conn
	natsJetStream nats.JetStream

	nopts *nats.Options

	lock sync.RWMutex
}

// NewConn .
func NewConn(opts ...Option) *conn {
	n := new(conn)
	n.opts = defaultOptions()
	for _, opt := range opts {
		opt(n.opts)
	}
	n.configura()
	return n
}

// configura .
func (n *conn) configura() {
	// default options for nats
	{
		nopts := nats.GetDefaultOptions()
		n.nopts = &nopts
	}

	// address
	n.nopts.Servers = n.addrs

	// TLSConfig
	if n.opts.TLSConfig != nil {
		n.nopts.Secure = true
		n.nopts.TLSConfig = n.opts.TLSConfig
	}
}

// Options .
func (n *conn) Options() *Options {
	return n.opts
}

// Address .
func (n *conn) Address() string {
	if n.natsConn != nil && n.natsConn.IsConnected() {
		return n.natsConn.ConnectedUrl()
	}
	return n.addrs[0]
}

// Connect .
func (n *conn) Connect() error {
	n.lock.Lock()

	status := nats.CLOSED
	if n.natsConn != nil {
		status = n.natsConn.Status()
	}

	switch status {
	case nats.CONNECTED, nats.CONNECTING, nats.RECONNECTING:
	default:
		// nats.Conn
		natsConn, err := n.nopts.Connect()
		if err != nil {
			n.lock.Unlock()
			return err
		}

		// nats.JetStream
		natsJetStream, err := natsConn.JetStream()
		if err != nil {
			n.lock.Unlock()
			return err
		}

		if _, err := natsJetStream.StreamInfo(n.opts.Stream); err == nats.ErrStreamNotFound {
			natsJetStream.AddStream(newStreamConfig(withStreamName(n.opts.Stream)))
		}

		n.natsConn = natsConn
		n.natsJetStream = natsJetStream
	}

	n.lock.Unlock()
	return nil
}

// Disconnect .
func (n *conn) Disconnect() error {
	n.lock.Lock()

	if n.natsConn != nil {
		n.natsConn.Flush()
		n.natsConn.Close()

		n.natsConn = nil
		n.natsJetStream = nil
	}

	n.lock.Unlock()
	return nil
}

// marshal .
func (n *conn) marshal(message *Message) ([]byte, error) {
	if message == nil || message.Data == nil {
		return nil, nats.ErrInvalidMsg
	}
	if message.Header == nil {
		message.Header = make(map[string]string)
	}
	message.Header["Content-Type"] = string(n.opts.Codec.String())
	return n.opts.Codec.Marshal(message)
}

// consumer .
func (n *conn) consumer(handler Handler) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		handler(&event{
			topic:   msg.Subject,
			reply:   msg.Reply,
			respond: msg.Respond,
			body:    msg.Data,
			codec:   n.opts.Codec,
		})
		msg.Ack()
	}
}

// Publish .
func (n *conn) Publish(topic string, message *Message, opts ...PublishOption) error {
	n.lock.RLock()

	if n.natsConn == nil {
		return nats.ErrInvalidConnection
	}

	body, err := n.marshal(message)
	if err != nil {
		return err
	}

	err = n.natsConn.Publish(topic, body)

	n.lock.RUnlock()
	return err
}

// Subscribe .
func (n *conn) Subscribe(topic string, handler Handler, opts ...SubscribeOption) error {
	n.lock.RLock()

	if n.natsConn == nil {
		return nats.ErrInvalidConnection
	}

	options := defaultSubscribeOptions()
	for _, opt := range opts {
		opt(options)
	}

	var err error

	if len(options.Queue) != 0 {
		_, err = n.natsConn.QueueSubscribe(topic, options.Queue, n.consumer(handler))
	} else {
		_, err = n.natsConn.Subscribe(topic, n.consumer(handler))
	}

	n.lock.RUnlock()
	return err
}

// JetStreamPublish .
func (n *conn) JetStreamPublish(topic string, message *Message, opts ...PublishOption) error {
	n.lock.RLock()

	if n.natsConn == nil {
		return nats.ErrInvalidConnection
	}

	body, err := n.marshal(message)
	if err != nil {
		return err
	}

	_, err = n.natsJetStream.Publish(
		topic,
		body,
		nats.ExpectStream(n.opts.Stream),
	)

	n.lock.RUnlock()
	return err
}

// JetStreamSubscribe .
func (n *conn) JetStreamSubscribe(topic string, handler Handler, opts ...SubscribeOption) error {
	n.lock.RLock()

	if n.natsConn == nil {
		return nats.ErrInvalidConnection
	}

	options := defaultSubscribeOptions()
	for _, opt := range opts {
		opt(options)
	}
	var err error

	if len(options.Queue) != 0 {
		_, err = n.natsJetStream.QueueSubscribe(
			topic,
			options.Queue,
			n.consumer(handler),
		)
	} else {
		_, err = n.natsJetStream.Subscribe(topic, n.consumer(handler))
	}

	n.lock.RUnlock()
	return err
}

type respond func(data []byte) error

// event .
type event struct {
	topic   string
	reply   string
	respond respond

	body   []byte
	header Header
	once   sync.Once
	codec  codec.Marshaler
}

func (e *event) Topic() string {
	return e.topic
}

// Reply .
func (e *event) Reply() string {
	return e.reply
}

// Body .
func (e *event) Body() []byte {
	return e.body
}

// Header .
func (e *event) Header() Header {
	e.once.Do(func() {
		message := new(Message)
		e.codec.Unmarshal(e.body, message)
		e.header = message.Header
	})
	return e.header
}

// Respond .
func (e *event) Respond(v interface{}) error {
	data, err := e.codec.Marshal(v)
	if err != nil {
		return err
	}
	return e.respond(data)
}

// Unmarshal .
func (e *event) Unmarshal(pointer interface{}) error {
	message := new(Message)
	message.Data = pointer

	return e.codec.Unmarshal(e.body, message)
}
