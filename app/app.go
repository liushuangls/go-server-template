package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/sourcegraph/conc"
	"golang.org/x/sync/errgroup"
)

type ServerStarter interface {
	OnServerStart(ctx context.Context) error
}

type ServerCloser interface {
	OnServerClose(ctx context.Context) error
}

type App struct {
	Options
}

func NewApp(opt Options) *App {
	app := &App{opt}
	return app
}

func (a *App) Run() error {
	signalCtx, signalCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer signalCancel()

	// start http server
	var g errgroup.Group
	httpSrv, err := a.Http.Run(&g)
	if err != nil {
		return err
	}
	// monitor http server error
	go func() {
		if err := g.Wait(); err != nil {
			slog.Error("http server err", slog.Any("err", err))
			signalCancel()
		}
	}()

	// start other servers
	if err := a.serverStart(signalCtx); err != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		_ = httpSrv.Shutdown(shutdownCtx)
		return err
	}

	fmt.Println("Server Started at", a.Http.Conf.App.Addr)

	select {
	case <-signalCtx.Done():
	}

	// end all services
	endCtx, endCancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer endCancel()

	var wg conc.WaitGroup

	// close http server
	wg.Go(func() {
		if err := httpSrv.Shutdown(endCtx); err != nil {
			slog.Error("Server Shutdown", "err", err)
		}
	})

	// close other servers
	a.serverClose(endCtx, &wg)

	if r := wg.WaitAndRecover(); r != nil {
		slog.Error("Server Shutdown", "wait err", r.String())
	}

	slog.Info("server exiting")
	return nil
}

func (a *App) serverStart(ctx context.Context) error {
	v := reflect.ValueOf(a.Options)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if (field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface) && field.IsNil() {
			continue
		}
		if starter, ok := field.Interface().(ServerStarter); ok {
			if err := starter.OnServerStart(ctx); err != nil {
				return fmt.Errorf("App.Run.serverStart type: %s, err: %w", reflect.TypeOf(starter).String(), err)
			}
		}
	}
	return nil
}

func (a *App) serverClose(ctx context.Context, wg *conc.WaitGroup) {
	v := reflect.ValueOf(a.Options)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if (field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface) && field.IsNil() {
			continue
		}
		if closer, ok := field.Interface().(ServerCloser); ok {
			localCloser := closer
			localType := reflect.TypeOf(localCloser).String()
			wg.Go(func() {
				if err := localCloser.OnServerClose(ctx); err != nil {
					slog.Error(
						"App.Run.serverClose",
						slog.Any("err", err),
						slog.String("type", localType),
					)
				}
			})
		}
	}
}
