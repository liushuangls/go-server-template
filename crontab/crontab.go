package crontab

import (
	"context"
	"log/slog"
	"reflect"
	"time"

	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/logger"
	"github.com/reugn/go-quartz/quartz"

	"github.com/liushuangls/go-server-template/pkg/ecode"
)

type Client struct {
	Options
	scheduler quartz.Scheduler
}

func NewClient(ctx context.Context, opt Options) (*Client, error) {
	sched, err := quartz.NewStdScheduler(quartz.WithLogger(logger.NewSlogLogger(ctx, slog.Default())))
	if err != nil {
		return nil, ecode.WithCaller(err)
	}
	return &Client{
		Options:   opt,
		scheduler: sched,
	}, nil
}

func (c *Client) Start() error {
	val := reflect.ValueOf(c)
	typ := val.Type()
	output1 := reflect.TypeOf((*quartz.JobDetail)(nil))
	output2 := reflect.TypeOf((*error)(nil)).Elem()
	for i := 0; i < typ.NumMethod(); i++ {
		method := val.Method(i)
		methodType := method.Type()
		if methodType.NumOut() == 2 && methodType.Out(0).AssignableTo(output1) && methodType.Out(1).Implements(output2) {
			results := method.Call(nil)
			if err, ok := results[1].Interface().(error); ok && err != nil {
				return err
			}
		}
	}

	c.scheduler.Start(context.Background())
	return nil
}

func (c *Client) Stop(ctx context.Context) {
	c.scheduler.Stop()
	c.scheduler.Wait(ctx)
}

func (c *Client) RegisterPrintTask() (*quartz.JobDetail, error) {
	detail := quartz.NewJobDetail(job.NewFunctionJob(func(ctx context.Context) (any, error) {
		slog.Info("crontab.RegisterPrintTask")
		return nil, nil
	}), quartz.NewJobKey("print_task"))
	err := c.scheduler.ScheduleJob(detail, quartz.NewSimpleTrigger(time.Minute))
	return detail, err
}
