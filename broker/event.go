package broker

import (
	"github.com/pkg/errors"

	"github.com/charlesbases/library/codec"
	"github.com/charlesbases/library/content"
)

// Event .
type Event interface {
	// Topic .
	Topic() string
	// Reply .
	Reply() string
	// Body return bytes of message
	Body() []byte
	// Unmarshal unmarshal message
	// 若消息格式为 json, 则反序列化 message.data; 若消息格式为 proto, 则反序列化 message
	Unmarshal(v interface{}) error
}

// event .
type event struct {
	topic string
	reply string
	body  []byte
	codec codec.Marshaler
}

// Topic .
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

// Unmarshal .
func (e *event) Unmarshal(v interface{}) error {
	switch e.codec.ContentType() {
	case content.Json:
		return e.codec.Unmarshal(e.body, &JsonMessage{Data: v})
	case content.Proto:
		return e.codec.Unmarshal(e.body, v)
	default:
		return errors.Errorf("unsupported of %s", e.codec.ContentType().String())
	}
}

// NewEvent .
func NewEvent(topic string, reply string, body []byte, codec codec.Marshaler) Event {
	return &event{topic: topic, reply: reply, body: body, codec: codec}
}
