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
	hooks []Hook
}

// New .
func New() *Lifecycle {
	lf := new(Lifecycle)
	lf.hooks = make([]Hook, 0, 4)
	return lf
}

// Append .
func (lf *Lifecycle) Append(h Hook) {
	lf.hooks = append(lf.hooks, h)
}

// Start .
func (lf *Lifecycle) Start(ctx context.Context) error {
	for _, h := range lf.hooks {
		if h.OnStart != nil {
			if err := h.OnStart(ctx); err != nil {
				logger.Errorf("%s start failed: %v", h.Name, err)
				return err
			}
		}
	}
	return nil
}

// Stop .
func (lf *Lifecycle) Stop(ctx context.Context) error {
	for _, h := range lf.hooks {
		if h.OnStop != nil {
			if err := h.OnStop(ctx); err != nil {
				logger.Errorf("%s stop failed: %v", h.Name, err)
				return err
			}
		}
	}
	return nil
}
