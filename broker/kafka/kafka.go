package kafka

import (
	"context"
	"fmt"
	"time"

	"library"
	"library/broker"
	"library/codec"

	"github.com/Shopify/sarama"
	"github.com/charlesbases/logger"
	"github.com/google/uuid"
)

// client .
type client struct {
	opts *broker.Options
	conf *sarama.Config

	addrs []string

	id string

	client   sarama.Client
	producer sarama.AsyncProducer

	actived bool
	closing chan struct{}
}

// Init .
func Init(id string, options ...broker.Option) broker.Broker {
	c := &client{id: id, closing: make(chan struct{})}
	c.configure(options...)
	return c
}

// configure .
func (c *client) configure(options ...broker.Option) {
	c.opts = new(broker.Options)
	for _, o := range options {
		o(c.opts)
	}

	c.addrs = c.opts.Addrs
	if len(c.addrs) == 0 {
		logger.Fatal(broker.ErrInvalidAddrs)
	}
	if c.opts.ReconnectTime == 0 {
		c.opts.ReconnectTime = broker.DefaultReconnectTime
	}

	c.conf = sarama.NewConfig()
	c.conf.ClientID = c.id

	c.conf.Consumer.Offsets.Initial = sarama.OffsetNewest
	c.conf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	c.conf.Consumer.Offsets.AutoCommit.Enable = true

	c.conf.Producer.Return.Errors = true
	c.conf.Producer.Return.Successes = false
	c.conf.Producer.RequiredAcks = sarama.WaitForAll
	c.conf.Producer.Partitioner = sarama.NewRandomPartitioner
}

// groupID .
func (c *client) groupID(topic string) string {
	return fmt.Sprintf("%s.%s", c.id, topic)
}

// newAsyncProducer .
func (c *client) newAsyncProducer() {
	// Producer async
	producer, err := sarama.NewAsyncProducerFromClient(c.client)
	if err != nil {
		logger.Fatalf("[Kafka] connect failed. NewProducer error: %v", err)
	}
	c.producer = producer
}

// newConsumerGroup .
func (c *client) newConsumerGroup(topic string) sarama.ConsumerGroup {
	consumer, err := sarama.NewConsumerGroupFromClient(c.groupID(topic), c.client)
	if err != nil {
		logger.Fatalf("[Kafka] connect failed. NewConsumerGroup error: %v", err)
	}
	logger.Infof(`[Kafka] subscribe["%s"]`, topic)
	return consumer
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
				logger.Errorf(`[Kafka] produce["%s"] failed. %v`, err.Msg.Topic, err.Err)
			}
		}
	}
}

// Address .
func (c *client) Address() []string {
	return c.addrs
}

// Publish 异步发布
func (c *client) Publish(topic string, v interface{}, opts ...broker.PublishOption) error {
	if c.actived {
		var opt = broker.DefaultPublishOptions()
		for _, o := range opts {
			o(opt)
		}

		message := &broker.Message{
			ID:          uuid.New(),
			Topic:       topic,
			Producer:    c.id,
			CreatedAt:   library.ParseTime2String(time.Now()),
			ContentType: opt.Codec.ContentType(),
			Data:        v,
		}
		bytes, err := opt.Codec.Marshal(message)
		if err != nil {
			logger.Errorf(`[Kafka] publish["%s"] failed. Message.Marshal error: %v`, topic, err)
			return err
		}

		c.producer.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(bytes),
		}

		logger.Debugf(`[Kafka] publish["%s"] %s`, topic, func() string {
			switch opt.Codec.ContentType() {
			case codec.ContentTypeProto:
				return "[ProtoMessage]"
			default:
				return string(bytes)
			}
		}())
	} else {
		logger.Warnf(`[Kafka] publish["%s"] failed. broker is disconnected`, topic)
	}
	return nil
}

// Subscribe 消息订阅
func (c *client) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) {
	if c.actived {
		var opt = broker.DefaultSubscribeOptions()
		for _, o := range opts {
			o(opt)
		}

		go func() {
			reconnect := time.NewTicker(c.opts.ReconnectTime)

			consumer := c.newConsumerGroup(topic)
			consumerGroupHandler := &consumerGroup{client: c, handler: handler, codec: opt.Codec}
			for {
				err := consumer.Consume(opt.Context, []string{topic}, consumerGroupHandler)
				select {
				case <-c.closing:
					consumer.Close()
					return
				default:
					if err != nil {
						logger.Errorf(`[Kafka] consume["%s"] failed. %v`, topic, err)
					}
					<-reconnect.C
				}
			}
		}()
	} else {
		logger.Warnf(`[Kafka] subscribe["%s"] failed. broker is disconnected`, topic)
	}
}

// OnStart .
func (c *client) OnStart(ctx context.Context) error {
	if !c.actived {
		// Client
		client, err := sarama.NewClient(c.addrs, c.conf)
		if err != nil {
			logger.Fatalf("[Kafka] connect failed. %v", err)
		}
		c.client = client

		c.newAsyncProducer()

		c.actived = true

		go c.daemon()
	}
	return nil
}

// OnStop .
func (c *client) OnStop(ctx context.Context) error {
	if c.actived {
		c.actived = false
		close(c.closing)
		c.client.Close()
	}
	return nil
}

// String .
func (c *client) String() string {
	return "Kafka"
}

// consumerGroup .
type consumerGroup struct {
	client  *client
	handler broker.Handler

	codec codec.Marshaler
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

			logger.Debugf(`[Kafka] consume["%s"] %s`, message.Topic, func() string {
				switch c.codec.ContentType() {
				case codec.ContentTypeProto:
					return "[ProtoMessage]"
				default:
					return string(message.Value)
				}
			}())

			go func() {
				err := c.handler(&event{
					topic: message.Topic,
					body:  message.Value,
					codec: c.codec,
				})

				if err != nil {
					logger.Errorf(`[Kafka] consume["%s"] failed: %s`, message.Topic, err.Error())
				} else {
					session.MarkMessage(message, "")
				}
			}()
		}
	}
}

// event .
type event struct {
	topic string
	body  []byte
	codec codec.Marshaler
}

// Topic .
func (e *event) Topic() string {
	return e.topic
}

// Body .
func (e *event) Body() []byte {
	return e.body
}

// Unmarshal .
func (e *event) Unmarshal(v interface{}) error {
	var message = &broker.Message{Data: v}
	return e.codec.Unmarshal(e.body, message)
}
