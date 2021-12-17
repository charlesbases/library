package nats

import (
	"fmt"
	"testing"
	"time"
)

var handler = func(event Event) interface{} {
	var data int
	event.Unmarshal(&data)
	return data
}

var consumerA = func(event Event) error {
	fmt.Printf("  A --> %s.%v\n", event.Topic(), handler(event))
	return nil
}

var consumerB = func(event Event) error {
	fmt.Printf("  B --> %s.%v\n", event.Topic(), handler(event))
	return nil
}

var consumerC = func(event Event) error {
	fmt.Printf("  C --> %s.%v\n", event.Topic(), handler(event))
	return nil
}

func Test(t *testing.T) {
	nats := NewConn()
	if err := nats.Connect(); err != nil {
		panic(err)
	}

	var (
		topicA = "normal" // normal
		topicB = "stream" // stream
	)

	go func() {
		<-time.NewTicker(time.Second * 3).C
		{
			// nats.Subscribe("*", consumerA)
			// nats.Subscribe("*", consumerB)
			// nats.Subscribe("*", consumerC)
			// nats.Subscribe(topicA, consumerA)
			// nats.Subscribe(topicA, consumerB)
			// nats.Subscribe(topicA, consumerC)
			// nats.Subscribe(topicA, consumerA, WithSubscribeQueue("queue"))
			// nats.Subscribe(topicA, consumerB, WithSubscribeQueue("queue"))
			// nats.Subscribe(topicA, consumerC, WithSubscribeQueue("queue"))
		}
		{
			// nats.JetStreamSubscribe("*", consumerA)
			// nats.JetStreamSubscribe("*", consumerB)
			// nats.JetStreamSubscribe("*", consumerC)
			nats.JetStreamSubscribe(topicB, consumerA)
			nats.JetStreamSubscribe(topicB, consumerB)
			nats.JetStreamSubscribe(topicB, consumerC)
		}
	}()

	var heartbeat int = 0
	for heartbeat < 10 {
		heartbeat++

		{
			err := nats.Publish(topicA, &Message{
				Data: heartbeat,
			})
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("publish --> %s.%d\n", topicA, heartbeat)
		}

		{
			err := nats.JetStreamPublish(topicB, &Message{
				Data: heartbeat,
			})
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("publish --> %s.%d\n", topicB, heartbeat)
		}

		<-time.NewTicker(time.Second * 1).C
	}
}
