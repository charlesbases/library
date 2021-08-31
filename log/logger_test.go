package log

import (
	"testing"
	"time"
)

func TestSeelog(t *testing.T) {
	UseSeelog(WithFilename("./log.log"))

	Trace("trace")
	Debug("debug")
	Info("info")
	Error("error")

	<-time.NewTicker(time.Second).C
}

func TestZap(t *testing.T) {
	UseZap()

	Trace("trace")
	Debug("debug")
	Info("info")
	Error("error")
}
