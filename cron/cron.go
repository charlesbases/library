package cron

import (
	"context"

	"github.com/charlesbases/logger"
	"github.com/robfig/cron/v3"

	"github.com/charlesbases/library/lifecycle"
)

type Stop func() context.Context

var background = func() context.Context {
	return context.Background()
}

type CronCmd interface {
	AddFunc(spec string, action func()) (cron.EntryID, error)
}

// clog .
type clog struct{}

// Info .
func (clog) Info(_ string, _ ...interface{}) {
}

// Error .
func (clog) Error(err error, _ string, _ ...interface{}) {
	logger.Named("crontab").Error(err)
}

// Run 添加并运行定时任务，并返回 jobs 的 Stop 函数。
func Run(fn func(cron CronCmd) error) (Stop, error) {
	c := cron.New(cron.WithLogger(new(clog)))
	if err := fn(c); err != nil {
		return background, err
	}
	c.Start()
	return c.Stop, nil
}

// Append 添加定时任务，并且随着 lifecycle.Lifecycle 运行或停止
func Append(lf *lifecycle.Lifecycle, fn func(cron CronCmd) error) error {
	c := cron.New(cron.WithLogger(new(clog)))
	if err := fn(c); err != nil {
		return err
	}

	lf.Append(&lifecycle.Hook{
		Name: "crontab",
		OnStart: func(ctx context.Context) error {
			c.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			c.Stop()
			return nil
		},
	})
	return nil
}
