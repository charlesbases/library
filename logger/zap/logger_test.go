package zap

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	go func() {
		{
			Debug(now())
			Info(now())
			Warn(now())
			Error(now())
		}
	}()

	go func() {
		{
			l := New(WithService("1"))
			l.Debug(now())
			l.Info(now())
			l.Warn(now())
			l.Error(now())
		}
	}()

	go func() {
		{
			l := New(WithService("2"))
			l.Debug(now())
			l.Info(now())
			l.Warn(now())
			l.Error(now())
		}
	}()

	<-time.After(time.Second * 3)
}

// now .
func now() string {
	return time.Now().Format(DefaultDateFormat)
}
