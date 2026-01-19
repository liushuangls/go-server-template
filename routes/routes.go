package routes

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/go-redis/redis_rate/v10"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/routes/common"
	"github.com/liushuangls/go-server-template/routes/middleware"
)

func NewEcho(conf *configs.Config, logger *slog.Logger) (*echo.Echo, error) {
	e := echo.New()

	e.Logger = logger
	e.HTTPErrorHandler = common.EchoErrorHandler

	cb, err := common.NewCustomBinder()
	if err != nil {
		return nil, err
	}
	e.Binder = cb

	e.Use(
		echoMiddleware.Recover(),
		echoMiddleware.RequestLogger(),
		echoMiddleware.CORS("*"),
	)

	return e, nil
}

type HttpEngine struct {
	Options
}

type Registrable interface {
	RegisterRoute(group *echo.Group)
}

func NewHttpEngine(opt Options) *HttpEngine {
	return &HttpEngine{opt}
}

func (h *HttpEngine) RegisterRoute() {
	g := h.Router.Group("")
	g.Use(
		middleware.RateLimitWithIP(h.Limiter, redis_rate.PerMinute(60), "total"),
	)

	v := reflect.ValueOf(h.Options)
	for i := 0; i < v.NumField(); i++ {
		if router, ok := v.Field(i).Interface().(Registrable); ok {
			router.RegisterRoute(g)
		}
	}

	printRoutes(h.Router)
}

func printRoutes(e *echo.Echo) {
	fmt.Println("==== Registered Routes ====")
	for _, r := range e.Router().Routes() {
		if r.Path == "/" || r.Path == "/*" {
			continue
		}
		// r.Name 是 handler 的函数名，视情况打印
		fmt.Printf("%-6s %-30s -> %s\n", r.Method, r.Path, r.Name)
	}
	fmt.Println("===========================")
}

func (h *HttpEngine) Run() (*http.Server, error) {
	h.RegisterRoute()

	srv := &http.Server{
		Addr:    h.Conf.App.Addr,
		Handler: h.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	return srv, nil
}
