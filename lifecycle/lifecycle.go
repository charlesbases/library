package lifecycle

import "github.com/charlesbases/logger"

// Hook .
type Hook interface {
	OnStart() error
	OnStop() error
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
func (lf *Lifecycle) Start() error {
	for _, hook := range lf.hooks {
		if err := hook.OnStart(); err != nil {
			logger.Errorf("%s start failed: %v", hook.String(), err)
			return err
		}
	}

	return nil
}

// Stop .
func (lf *Lifecycle) Stop() error {
	for _, hook := range lf.hooks {
		if err := hook.OnStop(); err != nil {
			logger.Errorf("%s stop failed: %v", hook.String(), err)
			return err
		}
	}

	return nil
}
