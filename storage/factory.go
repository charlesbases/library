package storage

import (
	"sync"
)

const defaultFactoryWorks = 100

// Factory .
type Factory struct {
	// Concurrency The number of goroutines
	Works int
	keych chan interface{}

	once sync.Once
	swg  sync.WaitGroup
}

// NewFactory .
func NewFactory(options ...func(f *Factory)) *Factory {
	var f = &Factory{
		Works: defaultFactoryWorks,
		keych: make(chan interface{}, 8),
		swg:   sync.WaitGroup{},
	}

	for _, opt := range options {
		opt(f)
	}

	return f
}

// Flowline .
func (f *Factory) Flowline(fn func(v interface{})) {
	f.swg.Add(f.Works)

	for i := 0; i < f.Works; i++ {
		go func() {
			defer f.swg.Done()

			for {
				select {
				case key, ok := <-f.keych:
					if ok {
						fn(key)
					} else {
						return
					}
				}
			}
		}()
	}
}

// Push .
func (f *Factory) Push(v interface{}) {
	f.keych <- v
}

// PushSlice .
func (f *Factory) PushSlice(keys []*string) {
	var length = len(keys)

	// 一个协程只处理十个 key
	left, right := 0, 0
	for left != length {
		if right = left + 10; right >= length {
			right = length
		}

		f.Push(keys[left:right])
		left = right
	}
}

// Closing .
func (f *Factory) Closing() {
	f.once.Do(func() {
		close(f.keych)
	})
}

// Wait .
func (f *Factory) Wait() {
	f.Closing()
	f.swg.Wait()
}
