package cron

import (
	"fmt"
	"testing"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/lifecycle"
)

func TestCron(t *testing.T) {
	_, _ = Run(func(cron CronCmd) error {
		cron.AddFunc("0 0 0 * * *", func() {
			fmt.Println(library.NowString())
		})
		return nil
	})

	lf := new(lifecycle.Lifecycle)
	_ = Append(lf, func(cron CronCmd) error {
		cron.AddFunc("0 0 0 * * *", func() {
			fmt.Println(library.NowString())
		})
		return nil
	})

	lf.Start()
	lf.Stop()
}
