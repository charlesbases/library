package nats

import (
	"context"
	"strings"
	"sync"

	"library/broker"
	"library/codec"
	"library/codec/json"

	"github.com/nats-io/nats.go"
)

// natsBroker .
type natsBroker struct {
	addrs []string
	opts  *broker.Options

	conn  *nats.Conn
	nopts *nats.Options

	drain  bool
	closed chan (error)

	once sync.Once
	lock sync.RWMutex
}

// NewBroker .
func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.NewOptions(
		broker.WithCodec(json.NewMarshaler()),
		broker.WithContext(context.Background()),
	)

	n := new(natsBroker)
	n.opts = options
	n.configura(opts...)
	return n
}

// configura .
func (n *natsBroker) configura(opts ...broker.Option) {
	for _, opt := range opts {
		opt(n.opts)
	}

	// default options for nats
	n.once.Do(func() {
		opt := nats.GetDefaultOptions()
		n.nopts = &opt
	})

	// address
	{
		for _, addr := range n.opts.Addrs {
			if strings.HasPrefix(addr, "nats://") {
				n.addrs = append(n.addrs, addr)
			} else {
				n.addrs = append(n.addrs, "nats://"+addr)
			}
		}

		if len(n.addrs) == 0 {
			n.addrs = []string{nats.DefaultURL}
		}

		n.nopts.Servers = n.addrs
	}

	// TLSConfig
	{
		if n.opts.TLSConfig != nil {
			n.nopts.Secure = true
			n.nopts.TLSConfig = n.opts.TLSConfig
		}
	}

	n.drain = true
	n.closed = make(chan error)

	n.nopts.ClosedCB = n.onClose
	n.nopts.AsyncErrorCB = n.onAsyncError
}

// onClose .
func (n *natsBroker) onClose(conn *nats.Conn) {
	n.closed <- nil
}

// onAsyncError .
func (n *natsBroker) onAsyncError(conn *nats.Conn, sub *nats.Subscription, err error) {
	if err == nats.ErrDrainTimeout {
		n.closed <- err
	}
}

func (n *natsBroker) Init(opts ...broker.Option) error {
	n.configura(opts...)
	return nil
}

func (n *natsBroker) Options() *broker.Options {
	return n.opts
}

func (n *natsBroker) Address() string {
	if n.conn != nil && n.conn.IsConnected() {
		return n.conn.ConnectedUrl()
	}
	if len(n.addrs) > 0 {
		return n.addrs[0]
	}
	return ""
}

func (n *natsBroker) Connect() error {
	n.lock.Lock()
	defer n.lock.Unlock()

	status := nats.CLOSED
	if n.conn != nil {
		status = n.conn.Status()
	}

	switch status {
	case nats.CONNECTED, nats.RECONNECTING, nats.CONNECTING:
		return nil
	default:
		if conn, err := n.nopts.Connect(); err != nil {
			return err
		} else {
			n.conn = conn
		}
		return nil
	}
}

func (n *natsBroker) Disconnect() error {
	n.lock.RLock()
	defer n.lock.RUnlock()

	if n.drain {
		n.conn.Drain()
		return <-n.closed
	}
	n.conn.Close()
	return nil
}

func (n *natsBroker) Publish(topic string, message *broker.Message, opts ...broker.PublishOption) error {
	if message.Header == nil {
		message.Header = broker.NewHander()
	}
	message.Header["Content-Type"] = string(n.opts.Codec.String())

	body, err := n.opts.Codec.Marshal(message)
	if err != nil {
		return err
	}

	n.lock.RLock()
	err = n.conn.Publish(topic, body)
	n.lock.RUnlock()

	return err
}

func (n *natsBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) error {
	if n.conn == nil {
		return nats.ErrInvalidConnection
	}

	options := broker.NewSubscribeOptions(opts...)

	cb := func(msg *nats.Msg) {
		handler(&publication{topic: msg.Subject, reply: msg.Reply, body: msg.Data, codec: n.opts.Codec})
	}

	var err error

	n.lock.RLock()
	if len(options.Queue) != 0 {
		_, err = n.conn.QueueSubscribe(topic, options.Queue, cb)
	} else {
		_, err = n.conn.Subscribe(topic, cb)
	}
	n.lock.RUnlock()

	return err
}

func (n *natsBroker) String() string {
	return "nats"
}

// publication .
type publication struct {
	topic string
	reply string

	body   []byte
	header broker.Header

	once  sync.Once
	codec codec.Marshaler
}

func (p *publication) Topic() string {
	return p.topic
}

// Reply .
func (p *publication) Reply() string {
	return p.reply
}

// Body .
func (p *publication) Body() []byte {
	return p.body
}

// Header .
func (p *publication) Header() broker.Header {
	p.once.Do(func() {
		message := new(broker.Message)
		p.codec.Unmarshal(p.body, message)
		p.header = message.Header
	})

	return p.header
}

// Unmarshal .
func (p *publication) Unmarshal(pointer interface{}) error {
	message := new(broker.Message)
	message.Data = pointer

	return p.codec.Unmarshal(p.body, message)
}
