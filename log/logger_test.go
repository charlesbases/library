package log

import (
	"testing"
)

func TestSeelog(t *testing.T) {
	UseSeelog(WithFilename("./log.log"))

	Trace("trace")
	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")
}

func TestZap(t *testing.T) {
	UseZap(WithService("test"))

	Trace("trace")
	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")
}
