package message

import "testing"

func TestEmail(t *testing.T) {

	t.Run("smtp.163.com", func(t *testing.T) {
		email := NewEmailDialer("smtp.163.com", 465, "username", "passwd")
		email.Send(nil)
	})

	t.Run("smtp.qq.com", func(t *testing.T) {
		email := NewEmailDialer("smtp.qq.com", 465, "username", "passwd")
		email.Send(nil)
	})

	t.Run("smtp.exmail.qq.com", func(t *testing.T) {
		email := NewEmailDialer("smtp.exmail.qq.com", 465, "username", "passwd")
		email.Send(nil)
	})
}
