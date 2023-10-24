package lifecycle

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func TestLifecycle(t *testing.T) {
	lf := new(Lifecycle)

	lf.Append(
		&Hook{
			Name: "a",
			OnStart: func(ctx context.Context) error {
				fmt.Println("program a start")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				fmt.Println("program a stop")
				return nil
			},
		},
		&Hook{
			Name: "b",
			OnStart: func(ctx context.Context) error {
				return errors.New("test start error")
			},
			OnStop: func(ctx context.Context) error {
				return errors.New("test stop error")
			},
		},
	)

	lf.Start()
	lf.Stop()
}
