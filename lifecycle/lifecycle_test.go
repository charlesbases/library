package lifecycle

import (
	"fmt"
	"testing"
)

func TestLifecycle(t *testing.T) {
	lf := New()

	lf.Append(new(ServerA), new(ServerB))

	lf.Start()
	lf.Stop()
}

// ServerA .
type ServerA struct{}

// OnStart .
func (s *ServerA) OnStart() error {
	fmt.Println("A OnStart")
	return nil
}

// OnStop .
func (s *ServerA) OnStop() error {
	fmt.Println("A OnStop")
	return nil
}

// String .
func (s *ServerA) String() string {
	return "A"
}

// ServerB .
type ServerB struct{}

// OnStart .
func (s *ServerB) OnStart() error {
	fmt.Println("B OnStart")
	return nil
}

// OnStop .
func (s *ServerB) OnStop() error {
	fmt.Println("B OnStop")
	return nil
}

// String .
func (s *ServerB) String() string {
	return "B"
}
