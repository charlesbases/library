package nats

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/broker"
	"github.com/charlesbases/library/content"
	"github.com/charlesbases/library/logger"
)

// client .
type client struct {
	id   string
	opts *broker.Options

	conn   *nats.Conn
	stream nats.JetStreamContext

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
		return fmt.Errorf(`connect to "%s" failed. %v`, c.opts.Address, err)
	}

	stream, err := conn.JetStream()
	if err != nil {
		return fmt.Errorf(`connect to "%s" failed. %v`, c.opts.Address, err)
	}

	c.conn = conn
	c.stream = stream
	return nil
}

// orCreateStream add stream if stream name not existed
func (c *client) orCreateStream(v string) error {
	var err error
	if _, err = c.stream.StreamInfo(v); err == nats.ErrStreamNotFound {
		_, err = c.stream.AddStream(&nats.StreamConfig{
			Name:      v,
			Subjects:  []string{v},
			MaxAge:    7 * 24 * time.Hour,
			Retention: nats.WorkQueuePolicy,
		})
	}
	return err
}

func (c *client) Publish(topic string, v interface{}, opts ...func(o *broker.PublishOptions)) error {
	if !c.actived {
		return fmt.Errorf(`[nats] publish["%s"] failed. connection not ready.`, topic)
	}

	if err := broker.CheckTopic(topic); err != nil {
		return err
	}

	var o = broker.ParsePublishOptions(opts...)

	err := c.orCreateStream(topic)
	if err != nil {
		logger.ErrorfWithContext(o.Context, `[nats] publish["%s"] failed. %s`, topic, err.Error())
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
		err = fmt.Errorf("unsupported content-type of %v.", o.Codec.ContentType().String())
	}

	if err != nil {
		logger.ErrorfWithContext(o.Context, `[nats] publish["%s"] failed. %s`, topic, err.Error())
		return err
	}

	// publish
	if ack, err := c.stream.PublishMsgAsync(&nats.Msg{
		Subject: topic,
		Reply:   "", // TODO
		Data:    data,
	}, nats.ExpectStream(topic)); err != nil {
		logger.ErrorfWithContext(o.Context, `[nats] publish["%s"] failed. %s`, topic, err.Error())
	} else {
		go func() {
			select {
			case <-ack.Ok():
				logger.DebugfWithContext(o.Context, `[nats] publish["%s"] >> %s`, topic, o.Codec.ShowMessage(data))
			case err := <-ack.Err():
				logger.ErrorfWithContext(o.Context, `[nats] publish["%s"] failed. %s`, topic, err.Error())
			case <-time.NewTimer(5 * time.Second).C:
				logger.ErrorfWithContext(o.Context, `[nats] publish["%s"] failed. publish timeout.`, topic)
			}
		}()
	}
	return nil
}

func (c *client) Subscribe(topic string, handler broker.Handler, opts ...func(o *broker.SubscribeOptions)) {
	if !c.actived {
		logger.Errorf(`[nats] subscribe["%s"] failed. connection not ready.`, topic)
		return
	}

	if err := broker.CheckTopic(topic); err != nil {
		logger.Errorf(`[nats] subscribe["%s"] failed. %s.`, topic, err.Error())
		return
	}

	var o = broker.ParseSubscribeOptions(opts...)

	_, err := c.stream.QueueSubscribe(topic, o.ConsumerModel(c.id, topic),
		func(msg *nats.Msg) {
			msg.Ack()

			logger.Debugf(`[nats] consume["%s"] << %s`, msg.Subject, o.Codec.ShowMessage(msg.Data))

			if err := handler(broker.NewEvent(msg.Subject, msg.Reply, msg.Data, o.Codec)); err != nil {
				logger.Errorf(`[nats] consume["%s"] failed: %s`, msg.Subject, err.Error())
			}
		},
		nats.Durable(c.id+"."+topic))
	if err != nil {
		logger.Error(`[nats] subscribe["%s"] failed. %s.`, topic, err.Error())
		return
	}
}

func (c *client) Close() {
	if c.actived {
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
