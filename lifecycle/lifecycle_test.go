package lifecycle

import (
	"context"
	"fmt"
	"testing"
)

func TestLifecycle(t *testing.T) {
	lf := New()

	var a = new(ServerA)
	var b = new(ServerB)

	lf.Append(Hook{
		Name:    a.String(),
		OnStart: a.OnStart,
		OnStop:  a.OnStop,
	})

	lf.Append(Hook{
		Name: b.String(),
		OnStart: func(ctx context.Context) error {
			return b.OnStart()
		},
	})

	lf.Start(context.Background())
	lf.Stop(context.Background())
}

// ServerA .
type ServerA struct{}

// OnStart .
func (s *ServerA) OnStart(ctx context.Context) error {
	fmt.Println("A OnStart")
	return nil
}

// OnStop .
func (s *ServerA) OnStop(ctx context.Context) error {
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
