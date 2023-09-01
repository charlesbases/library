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

// options .
type options struct {
	ctx context.Context
}

type option func(o *options)

// Context .
func Context(ctx context.Context) option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// opts .
func (lf *Lifecycle) opts(opts ...option) *options {
	var o = &options{ctx: context.Background()}
	for _, opt := range opts {
		opt(o)
	}
	return o
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
func (lf *Lifecycle) Start(opts ...option) error {
	var opt = lf.opts(opts...)

	for _, hook := range lf.hooks {
		if hook.OnStart != nil {
			if err := hook.OnStart(opt.ctx); err != nil {
				logger.Errorf("[%s] start failed: %s", hook.Name, err.Error())
				return err
			}
		}
	}
	return nil
}

// Stop .
func (lf *Lifecycle) Stop(opts ...option) error {
	var opt = lf.opts(opts...)

	for _, hook := range lf.hooks {
		if hook.OnStop != nil {
			if err := hook.OnStop(opt.ctx); err != nil {
				logger.Errorf("[%s] stop failed: %s", hook.Name, err.Error())
				return err
			}
		}
	}
	return nil
}
