package library

import (
	"fmt"
	"sync"
	"testing"
)

func TestSignal(t *testing.T) {
	var sw sync.WaitGroup
	sw.Add(2)

	go func() {
		{
			c1 := Shutdown()
			select {
			case <-c1:
				fmt.Println("c1")
				close(c1)
				sw.Done()
			}
		}
	}()

	go func() {
		{
			c2 := Shutdown()
			select {
			case <-c2:
				fmt.Println("c2")
				close(c2)
				sw.Done()
			}
		}
	}()

	sw.Wait()
}
