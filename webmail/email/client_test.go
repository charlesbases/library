package email

import (
	"testing"
)

func Test(t *testing.T) {
	client, err := NewClient("", 0, "", "")
	if err != nil {
		panic(err)
	}
	client.Send(new(Message))
}
