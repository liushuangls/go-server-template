package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sourcegraph/conc"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Options
}

func NewApp(opt Options) *App {
	app := &App{opt}
	return app
}

func (a *App) Run() error {
	var g errgroup.Group

	// start http server
	httpSrv, err := a.Http.Run(&g)
	if err != nil {
		return err
	}
	// start crontab
	if err := a.Cron.Start(); err != nil {
		return err
	}

	fmt.Println("Server Started at", a.Http.Conf.App.Addr)

	quit := make(chan os.Signal, 1)

	// 监听中断信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 监控服务错误
	go func() {
		if err := g.Wait(); err != nil {
			slog.Error("server run err", "err", err)
			select {
			case quit <- syscall.SIGTERM:
			default:
			}
		}
	}()

	<-quit

	// 停止服务
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var wg conc.WaitGroup
	wg.Go(func() {
		a.Cron.Stop(ctx)
	})
	wg.Go(func() {
		if err := httpSrv.Shutdown(ctx); err != nil {
			slog.Error("Server Shutdown", "err", err)
		}
	})
	if r := wg.WaitAndRecover(); r != nil {
		slog.Error("Server Shutdown", "wait err", r.String())
	}

	slog.Info("server exiting")
	return nil
}
