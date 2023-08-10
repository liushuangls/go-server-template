package routes

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"go.uber.org/zap"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/internal/routes/common"
	"github.com/liushuangls/go-server-template/internal/routes/middleware"
	v1 "github.com/liushuangls/go-server-template/internal/routes/v1"
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

func (h *HttpEngine) Run() (*http.Server, error) {
	common.SetRespLog(h.log)
	h.RegisterRoute()
	srv := &http.Server{
		Addr:    h.conf.App.Addr,
		Handler: h.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			h.log.Fatalf("listen: %s\n", err)
		}
	}()
	return srv, nil
}
