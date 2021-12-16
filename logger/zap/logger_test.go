package zap

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	var loop int = 1e4

	logger := New(WithService("zap")) // 1.768725642s

	var start = time.Now()
	for i := 0; i < loop; i++ {
		logger.Debug(now())
		logger.Info(now())
		logger.Warn(now())
		logger.Error(now())
	}
	fmt.Println(time.Since(start))

	<-time.After(time.Second * 1)
}

// now .
func now() string {
	return time.Now().Format(DefaultDateFormat)
}
