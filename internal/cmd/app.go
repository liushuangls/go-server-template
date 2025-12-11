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
)

type App struct {
	Options
}

func NewApp(opt Options) *App {
	app := &App{opt}
	return app
}

func (a *App) Run() error {
	// start http server
	httpSrv, err := a.Http.Run()
	if err != nil {
		return err
	}
	// start crontab
	if err := a.Cron.Start(); err != nil {
		return err
	}

	fmt.Println("Server Started at", a.Http.Conf.App.Addr)

	// 监控结束指令
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 停止服务
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var wg conc.WaitGroup
	wg.Go(a.Cron.Stop)
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
