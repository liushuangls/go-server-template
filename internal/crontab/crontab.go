package crontab

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron"
)

type Client struct {
	Options
	scheduler *gocron.Scheduler
}

func NewClient(opt Options) *Client {
	s := gocron.NewScheduler(time.UTC)
	return &Client{
		Options:   opt,
		scheduler: s,
	}
}

func (c *Client) StartAsync() error {
	var (
		errFmt = "crontab.StartAsync error:%s"
		tasks  []func() (*gocron.Job, error)
	)

	tasks = append(tasks, c.registerPrintTask)

	for _, task := range tasks {
		if _, err := task(); err != nil {
			return fmt.Errorf(errFmt, err)
		}
	}

	c.scheduler.StartAsync()
	return nil
}

func (c *Client) Stop() {
	c.scheduler.Stop()
}

func (c *Client) registerPrintTask() (*gocron.Job, error) {
	return c.scheduler.Every("2m").Do(func() {
		slog.Info("crontab")
	})
}
