package nats

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/charlesbases/logger"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/broker"
	"github.com/charlesbases/library/content"
)

// client .
type client struct {
	id   string
	opts *broker.Options

	conn *nats.Conn
	js   nats.JetStreamContext

	actived bool
}

// connect .
func (c *client) connect() error {
	conn, err := nats.Connect(c.opts.Address, func(o *nats.Options) error {
		o.Name = c.id
		o.ReconnectWait = c.opts.ReconnectWait
		o.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		return nil
	})
	if err != nil {
		return err
	}

	js, err := conn.JetStream()
	if err != nil {
		return err
	}

	c.conn = conn
	c.js = js
	return nil
}

// orCreateStream add js if js name not existed
func (c *client) orCreateStream(v string) error {
	var err error
	if _, err = c.js.StreamInfo(v); err == nats.ErrStreamNotFound {
		_, err = c.js.AddStream(&nats.StreamConfig{
			Name:      v,
			Subjects:  []string{v},
			MaxAge:    7 * 24 * time.Hour,
			Retention: nats.WorkQueuePolicy,
		})
	}
	return err
}

// publish .
func (c *client) publish(subject string, v interface{}, opts ...func(o *broker.PublishOptions)) error {
	if !c.actived {
		return broker.ErrNotReady
	}

	if err := broker.CheckSubject(subject); err != nil {
		return err
	}

	var o = broker.ParsePublishOptions(opts...)

	err := c.orCreateStream(subject)
	if err != nil {
		return err
	}

	var data []byte
	switch o.Codec.ContentType() {
	case content.Json:
		data, err = o.Codec.Marshal(&broker.JsonMessage{
			ID:        uuid.NewString(),
			Producer:  c.id,
			CreatedAt: library.NowString(),
			Data:      v,
		})
	case content.Proto:
		data, err = o.Codec.Marshal(v)
	default:
		err = errors.Errorf("unsupported content-type of %v.", o.Codec.ContentType().String())
	}

	if err != nil {
		return err
	}

	// publish
	if ack, err := c.js.PublishMsgAsync(&nats.Msg{
		Subject: subject,
		Reply:   subject,
		Data:    data,
	}, nats.ExpectStream(subject)); err != nil {
		return err
	} else {
		go func() {
			select {
			case <-ack.Ok():
				logger.WithContext(o.Context).Debugf(`[nats]: publish["%s"]: %s`, subject, o.Codec.RawMessage(data))
			case err := <-ack.Err():
				logger.WithContext(o.Context).Errorf(`[nats]: publish["%s"]: %v`, subject, err)
			case <-time.NewTimer(o.Timeout).C:
				logger.WithContext(o.Context).Errorf(`[nats]: publish["%s"]: publish timeout`, subject)
			}
		}()
	}
	return nil
}

// Publish .
func (c *client) Publish(subject string, v interface{}, opts ...func(o *broker.PublishOptions)) error {
	return errors.Wrapf(c.publish(subject, v, opts...), `[nats]: publish["%s"]`, subject)
}

// subscribe .
func (c *client) subscribe(subject string, handler broker.Handler, opts ...func(o *broker.SubscribeOptions)) error {
	if !c.actived {
		return broker.ErrNotReady
	}

	if err := broker.CheckSubject(subject); err != nil {
		return err
	}

	logger.Debugf(`[nats]: subscribe["%s"]`, subject)

	var o = broker.ParseSubscribeOptions(opts...)
	_, err := c.js.QueueSubscribe(subject, o.ConsumerModel(c.id, subject),
		func(msg *nats.Msg) {
			msg.Ack()

			logger.Debugf(`[nats]: consume["%s"]: %s`, msg.Subject, o.Codec.RawMessage(msg.Data))
			if err := handler(broker.NewEvent(msg.Subject, msg.Reply, msg.Data, o.Codec)); err != nil {
				logger.Errorf(`[nats]: consume["%s"]: %v`, msg.Subject, err)
			}
		},
		nats.Durable(strings.Join([]string{c.id, subject}, ".")))
	return err
}

// Subscribe .
func (c *client) Subscribe(subject string, handler broker.Handler, opts ...func(o *broker.SubscribeOptions)) error {
	return errors.Wrapf(c.subscribe(subject, handler, opts...), `[nats]: subscribe["%s"]`, subject)
}

// Close .
func (c *client) Close() {
	if c.actived {
		c.actived = false

		c.conn.Flush()
		c.conn.Close()
	}
}

// NewClient .
func NewClient(id string, opts ...func(o *broker.Options)) (broker.Client, error) {
	c := &client{id: id, opts: broker.ParseOptions(opts...)}
	if len(c.opts.Address) == 0 {
		return nil, broker.ErrInvalidAddrs
	}
	return c, c.connect()
}
