package lifecycle

import (
	"context"

	"github.com/charlesbases/logger"
)

// Hook .
type Hook interface {
	OnStart(ctx context.Context) error
	OnStop(ctx context.Context) error
	String() string
}

// Lifecycle .
type Lifecycle struct {
	hooks []Hook
}

// New .
func New() *Lifecycle {
	lf := new(Lifecycle)
	lf.hooks = make([]Hook, 0, 4)
	return lf
}

// Append .
func (lf *Lifecycle) Append(hooks ...Hook) {
	if lf.hooks != nil {
		lf.hooks = append(lf.hooks, hooks...)
	} else {
		lf.hooks = hooks
	}
}

// Start .
func (lf *Lifecycle) Start(ctx context.Context) error {
	for _, hook := range lf.hooks {
		if err := hook.OnStart(ctx); err != nil {
			logger.Errorf("%s start failed: %v", hook.String(), err)
			return err
		}
	}

	return nil
}

// Stop .
func (lf *Lifecycle) Stop(ctx context.Context) error {
	for _, hook := range lf.hooks {
		if err := hook.OnStop(ctx); err != nil {
			logger.Errorf("%s stop failed: %v", hook.String(), err)
			return err
		}
	}

	return nil
}
