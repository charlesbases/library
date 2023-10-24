package kafka

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/charlesbases/logger"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/broker"
)

func Test(t *testing.T) {
	c, err := NewClient("test."+uuid.NewString(), func(o *broker.Options) {
		o.Address = "10.75.2.8:32509"
		o.Version = "3.3.1"
	})
	if err != nil {
		logger.Fatal(err)
	}

	var topic = "ticker"

	// Subscribe
	c.Subscribe(topic, func(event broker.Event) error {
		var timestr string
		if err := event.Unmarshal(&timestr); err != nil {
			return err
		}
		fmt.Println(timestr)
		return nil
	})

	// Publish
	go func() {
		for {
			c.Publish(topic, library.NowString(), func(o *broker.PublishOptions) {
				o.Context = context.WithValue(context.Background(), library.HeaderTraceID, uuid.NewString())
			})

			<-time.NewTimer(time.Second * 10).C
		}
	}()

	<-time.NewTimer(time.Minute).C
}
