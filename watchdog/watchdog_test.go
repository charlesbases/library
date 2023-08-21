package watchdog

import (
	"testing"
	"time"
)

func TestWatchDog(t *testing.T) {
	// os.Setenv(envMaxHeapMem, strconv.Itoa(1<<30))
	Memory()

	go func() {
		var v = make([]int, 0)

		for {
			v = append(v, 1)
		}
	}()

	<-time.NewTimer(time.Hour).C
}
