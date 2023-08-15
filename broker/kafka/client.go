package kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/broker"
	"github.com/charlesbases/library/content"
	"github.com/charlesbases/library/logger"
)

// consumerGroup .
type consumerGroup struct {
	client *client
	opts   *broker.SubscribeOptions

	h broker.Handler
}

func (c *consumerGroup) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroup) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-c.client.closing:
			return nil
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			session.MarkMessage(message, "")

			logger.Debugf(`[kafka] consume["%s"] << %s`, message.Topic, c.opts.Codec.ShowMessage(message.Value))

			go func() {
				if err := c.h(broker.NewEvent(message.Topic, message.Topic, message.Value, c.opts.Codec)); err != nil {
					logger.Errorf(`[kafka] consume["%s"] failed: %s`, message.Topic, err.Error())
				}
			}()
		}
	}
}

// client .
type client struct {
	id   string
	opts *broker.Options
	conf *sarama.Config

	client   sarama.Client
	producer sarama.AsyncProducer

	actived bool
	closing chan struct{}
}

// version .
func (c *client) version(ver *sarama.KafkaVersion) (err error) {
	*ver, err = sarama.ParseKafkaVersion(c.opts.Version)
	return err
}

func (c *client) connect() (err error) {
	// client
	c.client, err = sarama.NewClient([]string{c.opts.Address}, c.conf)
	if err != nil {
		return fmt.Errorf("connect failed. %v", err)
	}

	// producer
	c.producer, err = sarama.NewAsyncProducerFromClient(c.client)
	if err != nil {
		return fmt.Errorf("connect failed. NewProducer error: %v", err)
	}

	c.actived = true

	go c.daemon()
	return nil
}

// daemon .
func (c *client) daemon() {
	for {
		select {
		case <-c.closing:
			c.producer.Close()
			return
		case err, ok := <-c.producer.Errors():
			if ok {
				logger.Errorf(`[kafka] produce["%s"] failed. %v`, err.Msg.Topic, err.Err)
			}
		}
	}
}

func (c *client) Publish(topic string, v interface{}, opts ...func(o *broker.PublishOptions)) error {
	if !c.actived {
		return fmt.Errorf(`[kafka] publish["%s"] failed. connection not ready.`, topic)
	}

	if err := broker.CheckTopic(topic); err != nil {
		return err
	}

	var o = broker.ParsePublishOptions(opts...)

	var data []byte
	var err error
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
		err = fmt.Errorf("unsupported content-type of %s.", o.Codec.ContentType().String())
	}

	if err != nil {
		logger.ErrorfWithContext(o.Context, `[kafka] publish["%s"] failed. Marshal error: %s`, topic, err.Error())
		return err
	}

	c.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}

	logger.DebugfWithContext(o.Context, `[kafka] publish["%s"] >> %s`, topic, o.Codec.ShowMessage(data))
	return nil
}

func (c *client) Subscribe(topic string, handler broker.Handler, opts ...func(o *broker.SubscribeOptions)) {
	if !c.actived {
		logger.Errorf(`[kafka] subscribe["%s"] failed. connection not ready.`, topic)
		return
	}

	if err := broker.CheckTopic(topic); err != nil {
		logger.Errorf(`[kafka] subscribe["%s"] failed. %s`, err.Error())
		return
	}

	var o = broker.ParseSubscribeOptions(opts...)

	go func() {
		t := time.NewTicker(c.opts.ReconnectWait)

		consumer, err := sarama.NewConsumerGroupFromClient(o.ConsumerModel(c.id, topic), c.client)
		if err != nil {
			logger.Errorf(`[kafka] subscribe["%s"] failed. %s.`, err.Error())
			return
		}
		logger.Infof(`[kafka] subscribe["%s"]`, topic)

		consumerGroupHandler := &consumerGroup{client: c, h: handler, opts: o}
		for {
			err := consumer.Consume(o.Context, []string{topic}, consumerGroupHandler)
			select {
			case <-c.closing:
				consumer.Close()
				return
			default:
				if err != nil {
					logger.Errorf(`[kafka] consume["%s"] failed. %s`, topic, err.Error())
				}
				<-t.C
			}
		}
	}()
}

func (c *client) Close() {
	if c.actived {
		c.actived = false
		close(c.closing)
		c.client.Close()
	}
}

// NewClient .
func NewClient(id string, opts ...func(o *broker.Options)) (broker.Client, error) {
	c := &client{id: id, opts: broker.ParseOptions(opts...), closing: make(chan struct{})}

	if len(c.opts.Address) == 0 {
		return nil, broker.ErrInvalidAddrs
	}

	c.conf = sarama.NewConfig()
	if err := c.version(&c.conf.Version); err != nil {
		return nil, err
	}

	c.conf.ClientID = c.id

	c.conf.Consumer.Offsets.Initial = sarama.OffsetNewest
	c.conf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	c.conf.Consumer.Offsets.AutoCommit.Enable = true

	c.conf.Producer.Return.Errors = true
	c.conf.Producer.Return.Successes = false
	c.conf.Producer.RequiredAcks = sarama.WaitForAll
	c.conf.Producer.Partitioner = sarama.NewRandomPartitioner

	return c, c.connect()
}
