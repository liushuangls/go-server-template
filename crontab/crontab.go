package crontab

import (
	"log/slog"
	"reflect"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Client struct {
	Options
	scheduler gocron.Scheduler
}

func NewClient(opt Options) (*Client, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	return &Client{
		Options:   opt,
		scheduler: s,
	}, nil
}

func (c *Client) Start() error {
	val := reflect.ValueOf(c)
	typ := val.Type()
	taskType := reflect.TypeOf((*gocron.Job)(nil)).Elem()
	errType := reflect.TypeOf((*error)(nil)).Elem()
	for i := 0; i < typ.NumMethod(); i++ {
		method := val.Method(i)
		methodType := method.Type()
		if methodType.NumOut() == 2 && methodType.Out(0).Implements(taskType) && methodType.Out(1).Implements(errType) {
			results := method.Call(nil)
			if len(results) == 2 {
				if err, ok := results[1].Interface().(error); ok && err != nil {
					return err
				}
			}
		}
	}

	c.scheduler.Start()
	return nil
}

func (c *Client) Stop() {
	_ = c.scheduler.Shutdown()
}

func (c *Client) RegisterPrintTask() (gocron.Job, error) {
	return c.scheduler.NewJob(
		gocron.DurationJob(time.Minute*2),
		gocron.NewTask(func() {
			slog.Info("crontab")
		}),
	)
}
