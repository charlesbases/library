package sonyflake

import (
	"testing"

	"github.com/google/uuid"
)

func Benchmark(b *testing.B) {
	var bench = func(f func()) {
		b.ResetTimer()
		f()
		b.StopTimer()
	}

	b.Run("uuid", func(b *testing.B) {
		bench(func() {
			for i := 0; i < 10000; i++ {
				uuid.New().ID()
			}
		})
	})

	b.Run("sonyflake", func(b *testing.B) {
		bench(func() {
			for i := 0; i < 10000; i++ {
				NextID()
			}
		})
	})
}
