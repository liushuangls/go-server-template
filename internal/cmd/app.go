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

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/pkg/jwt"
	"github.com/liushuangls/go-server-template/pkg/xslog"
)

type App struct {
	Options
}

func NewApp(opt Options) *App {
	app := &App{opt}
	app.setDefaultSlog()
	return app
}

func (a *App) setDefaultSlog() {
	var extraWriters []xslog.ExtraWriter

	if a.Config.IsDebugMode() {
		a.Config.Log.Level = slog.LevelDebug
		extraWriters = append(extraWriters, xslog.ExtraWriter{
			Writer: os.Stdout,
			Level:  slog.LevelDebug,
		})
	}

	a.Config.Log.ExtraWriters = extraWriters
	fileLogger := xslog.NewFileSlog(&a.Config.Log)
	slog.SetDefault(fileLogger)
}

func NewJwt(conf *configs.Config) (*jwt.JWT, error) {
	return jwt.NewJWT(&jwt.Config{Issuer: conf.Jwt.Issuer, Secret: conf.Jwt.Secret})
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
