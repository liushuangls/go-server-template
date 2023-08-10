package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sourcegraph/conc"
	"go.uber.org/zap"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/pkg/jwt"
	"github.com/liushuangls/go-server-template/pkg/logger"
)

type App struct {
	Options
}

func NewApp(opt Options) *App {
	return &App{opt}
}

func NewLogger(conf *configs.Config) *zap.SugaredLogger {
	return logger.NewLogger(&conf.Log)
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
	if err := a.Cron.StartAsync(); err != nil {
		return err
	}

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
			a.Log.Errorw("Server Shutdown", "err", err)
		}
	})
	if r := wg.WaitAndRecover(); r != nil {
		a.Log.Errorw("Server Shutdown", "wait err", r.String())
	}

	a.Log.Infof("server exiting")
	return nil
}
