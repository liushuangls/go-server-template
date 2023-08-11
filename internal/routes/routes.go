package routes

import (
	"errors"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/internal/routes/middleware"
)

func NewEngine(conf *configs.Config) *gin.Engine {
	if conf.IsReleaseMode() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.Cors(true),
	)
	return r
}

type HttpEngine struct {
	Options
}

type Registrable interface {
	RegisterRoute(*gin.RouterGroup)
}

func NewHttpEngine(opt Options) *HttpEngine {
	return &HttpEngine{opt}
}

func (h *HttpEngine) RegisterRoute() {
	r := h.Router.Group("")
	r.Use(middleware.RateLimitWithIP(h.Limiter, redis_rate.PerMinute(60), "total"))

	v := reflect.ValueOf(h.Options)
	for i := 0; i < v.NumField(); i++ {
		if router, ok := v.Field(i).Interface().(Registrable); ok {
			router.RegisterRoute(r)
		}
	}
}

func (h *HttpEngine) Run() (*http.Server, error) {
	h.RegisterRoute()
	srv := &http.Server{
		Addr:    h.Conf.App.Addr,
		Handler: h.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return srv, nil
}
