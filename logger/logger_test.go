package logger

import (
	"testing"
)

func TestSeelog(t *testing.T) {
	NewSeelog()

	Trace("trace")
	Debug("debug")
	Info("info")
	Error("error")
}
