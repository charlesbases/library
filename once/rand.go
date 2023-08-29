package once

import (
	"math/rand"
	"sync"
)

/*
防止多个代码块中重复调用相关函数
*/

var randSeed sync.Once

// RandSeed .
func RandSeed(seed int64) {
	randSeed.Do(func() {
		rand.Seed(seed)
	})
}
