package log

import (
	"fmt"
	"testing"
	"time"
)

func TestSeelog(t *testing.T) {
	log := NewSeelog(
		WithFilename("./log.log"),
		WithService("test"),
	)

	start := time.Now()
	for i := 0; i < 10000; i++ {
		log.Trace("t", "r", "c", i)
		log.Debug("d", "b", "g", i)
		log.Info("i", "n", "f", i)
		log.Warn("w", "r", "n", i)
		log.Error("e", "r", "r", i)
	}
	fmt.Println("耗时: ", time.Now().Sub(start)) // 2.33s
}

func TestZap(t *testing.T) {
	log := NewZap(WithService("test"))

	start := time.Now()
	for i := 0; i < 10000; i++ {
		log.Trace("t", "r", "c", i)
		log.Debug("d", "b", "g", i)
		log.Info("i", "n", "f", i)
		log.Warn("w", "r", "n", i)
		log.Error("e", "r", "r", i)
	}
	fmt.Println("耗时: ", time.Now().Sub(start)) // 1.60s
}
