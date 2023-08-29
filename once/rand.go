package once

import (
	"math/rand"
	"sync"
)

var randSeed sync.Once

// RandSeed 防止多个代码块中重复调用 rand.Seed()
func RandSeed(seed int64) {
	randSeed.Do(func() {
		rand.Seed(seed)
	})
}
