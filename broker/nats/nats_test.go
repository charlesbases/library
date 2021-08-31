package nats

import (
	"fmt"
	"testing"
	"time"

	"library/broker"
)

var handler = func(event broker.Event) {
	date := new(Date)
	event.Unmarshal(date)

	fmt.Println(event.Header())
	fmt.Println(fmt.Sprintf(`%+v`, date))
}

// Date .
type Date struct {
	Format string
}

// NewDate .
func NewDate() *Date {
	return &Date{Format: time.Now().Format("2006-01-02 15:04:05")}
}

func Test(t *testing.T) {
	nats := NewBroker()
	nats.Connect()

	nats.Subscribe("date", handler)

	for {
		select {
		case <-time.NewTicker(time.Second * 3).C:
			err := nats.Publish("date", &broker.Message{
				Data: NewDate(),
			})
			if err != nil {
				fmt.Println("nats: publish error: ", err)
				break
			}
		}
	}
}
