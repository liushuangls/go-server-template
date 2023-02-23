package routes

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	configs "github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/internal/routes/common"
	"github.com/liushuangls/go-server-template/internal/routes/middleware"
	v1 "github.com/liushuangls/go-server-template/internal/routes/v1"
	"go.uber.org/zap"
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
	router  *gin.Engine
	conf    *configs.Config
	log     *zap.SugaredLogger
	limiter *redis_rate.Limiter

	user *v1.UserRoute
}

func NewHttpEngine(router *gin.Engine, conf *configs.Config, log *zap.SugaredLogger, limiter *redis_rate.Limiter, user *v1.UserRoute) *HttpEngine {
	return &HttpEngine{router: router, conf: conf, log: log, limiter: limiter, user: user}
}

func (h *HttpEngine) RegisterRoute() {
	r := h.router.Group("")
	r.Use(middleware.RateLimitWithIP(h.limiter, redis_rate.PerMinute(60), "total"))

	h.user.RegisterRoute(r)
}

func (h *HttpEngine) Run() error {
	common.SetRespLog(h.log)
	h.RegisterRoute()
	srv := &http.Server{
		Addr:    h.conf.App.Addr,
		Handler: h.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			h.log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		h.log.Fatal("Server Shutdown:", err)
	}
	h.log.Infof("server exiting")
	return nil
}
