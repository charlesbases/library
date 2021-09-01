package log

import (
	"fmt"
	"testing"
	"time"
)

func TestSeelog(t *testing.T) {
	UseSeelog(
		WithFilename("./log.log"),
		WithService("test"),
	)

	start := time.Now()
	for i := 0; i < 10000; i++ {
		Trace("t", "r", "c", i)
		Debug("d", "b", "g", i)
		Info("i", "n", "f", i)
		Warn("w", "r", "n", i)
		Error("e", "r", "r", i)
	}
	fmt.Println("耗时: ", time.Now().Sub(start)) // 2.33s
}

func TestZap(t *testing.T) {
	UseZap(WithService("test"))

	start := time.Now()
	for i := 0; i < 10000; i++ {
		Trace("t", "r", "c", i)
		Debug("d", "b", "g", i)
		Info("i", "n", "f", i)
		Warn("w", "r", "n", i)
		Error("e", "r", "r", i)
	}
	fmt.Println("耗时: ", time.Now().Sub(start)) // 1.60s
}
