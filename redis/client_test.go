package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/charlesbases/library/logger"

	"github.com/charlesbases/library"
)

func TestClient(t *testing.T) {
	ctx := context.WithValue(context.Background(), library.HeaderTraceID, uuid.NewString())

	Init(func(o *Options) {
		o.Addrs = []string{"10.63.2.46:6379"}
		o.Password = "admin123456.."
		o.Context = ctx
	})

	var fnGet = func(key keyword) {
		output := Client().Get(key)
		if output.err != nil {
			// logger.Fatal(err)
		}

		var v string
		if output.Unmarshal(&v) == nil {
			fmt.Println("value :", v)
			fmt.Println("ttl   :", output.TTL())
			fmt.Println("expire:", output.Expiry())
		}
	}

	// var fnDel = func(key string) {
	// 	r.Del(key)
	//
	// 	fnGet(key)
	// }

	keyprefix := KeyPrefix("t_")
	key := keyprefix("time")

	// Set
	{
		Client().Set(key, library.NowString(), func(o *SetOptions) {
			o.TTL = 3 * time.Second
			o.Context = ctx
		})

		fnGet(key)
	}

	// Del
	// {
	// 	fmt.Println()
	// 	fnDel(key)
	// }

	// Expire
	// {
	// 	fmt.Println()
	// 	r.Expire(key, func(o *ExpireOptions) {
	// 		o.TTL = 6 * time.Second
	// 	})
	//
	// 	fnGet(key)
	// }

	// Mutex
	{
		rm := Client().Mutex(key, func(o *MutexOptions) {
			o.Context = ctx
		})
		go func() {
			rm.Lock()
			logger.Debug("lock 1")
			rm.Unlock()
		}()

		go func() {
			rm.Lock()
			logger.Debug("lock 2")
			rm.Unlock()
		}()

		go func() {
			rm.Lock()
			logger.Debug("lock 3")
			rm.Unlock()
		}()

		go func() {
			rm.Lock()
			logger.Debug("lock 4")
			rm.Unlock()
		}()

		go func() {
			rm.Lock()
			logger.Debug("lock 5")
			rm.Unlock()
		}()

		<-time.NewTimer(time.Minute).C
	}
}

func TestCluster(t *testing.T) {

}
