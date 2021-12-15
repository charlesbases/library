package zap

import "testing"

func TestFileWriter(t *testing.T) {
	logger := New(WithService("FileWriter"))
	logger.Debug(now())
	logger.Info(now())
	logger.Warn(now())
	logger.Error(now())
}
