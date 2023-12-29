package hfwctx

import (
	"testing"

	"github.com/charlesbases/logger"
	"github.com/gin-gonic/gin"
)

func Test(t *testing.T) {
	_, err := Get("http://0.0.0.0:8080/api",
		Param("key", "val1", "val2", "var3", "var4", "var5"),
		Header("Content-A", "application/json"),
		Header("Content-B", "application/json"),
		Header("Content-C", "application/json"),
		Header("Content-D", "application/json"),
		Header("Content-E", "application/json"),
		Context(new(gin.Context)))
	if err != nil {
		logger.Error(err)
	}

	Get("http://0.0.0.0:8080/api",
		Param("key", "val1", "val2", "var3", "var4", "var5"),
		Header("Content-A", "application/json"),
		Header("Content-B", "application/json"),
		Header("Content-C", "application/json"),
		Header("Content-D", "application/json"),
		Header("Content-E", "application/json"),
		Context(new(gin.Context)))
}

// Benchmark-16    	     502	   2257775 ns/op	   20048 B/op	     132 allocs/op
func Benchmark(b *testing.B) {
	var bench = func(f func()) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f()
		}
		b.StopTimer()
	}

	bench(func() {
		Get("http://0.0.0.0:8080/api", Param("key", "val1", "val2"), Header("Content-type", "application/json"), Context(new(gin.Context)))
	})
}
