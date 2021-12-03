package nats

import (
	"fmt"
	"testing"
	"time"

	"library/broker"
)

var handler = func(event broker.Event) error {
	date := new(Date)
	event.Unmarshal(date)

	fmt.Println("receive -->", date.Format)
	return nil
}

// Date .
type Date struct {
	Format string
}

// NewDate .
func NewDate() *Date {
	return &Date{Format: time.Now().Format("2006-01-02 15:04:05.000")}
}

func Test(t *testing.T) {
	nats := NewBroker()
	nats.Connect()

	nats.Subscribe("date", handler)

	for {
		select {
		case <-time.NewTicker(time.Second * 3).C:
			date := NewDate()

			err := nats.Publish("date", &broker.Message{
				Data: date,
			})
			if err != nil {
				fmt.Println("nats: publish error: ", err)
				break
			}
			fmt.Println("publish -->", date.Format)
		}
	}
}
