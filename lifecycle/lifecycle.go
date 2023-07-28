package lifecycle

import (
	"context"

	"github.com/charlesbases/logger"
)

// Hook .
type Hook struct {
	Name string

	OnStart func(ctx context.Context) error
	OnStop  func(ctx context.Context) error
}

// Lifecycle .
type Lifecycle struct {
	hooks []*Hook
}

// Append .
func (lf *Lifecycle) Append(hooks ...*Hook) {
	if len(lf.hooks) != 0 {
		lf.hooks = append(lf.hooks, hooks...)
	} else {
		lf.hooks = hooks
	}
}

// Start .
func (lf *Lifecycle) Start(ctx context.Context) error {
	for _, hook := range lf.hooks {
		if hook.OnStart != nil {
			if err := hook.OnStart(ctx); err != nil {
				logger.Errorf("[%s] start failed: %s", hook.Name, err.Error())
				return err
			}
		}
	}
	return nil
}

// Stop .
func (lf *Lifecycle) Stop(ctx context.Context) error {
	for _, hook := range lf.hooks {
		if hook.OnStop != nil {
			if err := hook.OnStop(ctx); err != nil {
				logger.Errorf("[%s] stop failed: %s", hook.Name, err.Error())
				return err
			}
		}
	}
	return nil
}
