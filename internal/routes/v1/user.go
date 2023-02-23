package v1

import (
	"github.com/bsm/redislock"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/liushuangls/go-server-template/internal/dto/request"
	"github.com/liushuangls/go-server-template/internal/routes/common"
	"github.com/liushuangls/go-server-template/internal/routes/middleware"
	"github.com/liushuangls/go-server-template/internal/service"
	"go.uber.org/zap"
)

type UserRoute struct {
	log         *zap.SugaredLogger
	limiter     *redis_rate.Limiter
	locker      *redislock.Client
	userService *service.UserService
}

func NewUserRoute(log *zap.SugaredLogger, limiter *redis_rate.Limiter, locker *redislock.Client, userService *service.UserService) *UserRoute {
	return &UserRoute{log: log, limiter: limiter, locker: locker, userService: userService}
}

func (u *UserRoute) RegisterRoute(router *gin.RouterGroup) {
	user := router.Group("/v1/user")
	{
		user.POST(
			"/login",
			middleware.RateLimitWithIP(u.limiter, redis_rate.PerMinute(5), "login"),
			u.loginWithEmail,
		)
	}
}

func (u *UserRoute) loginWithEmail(c *gin.Context) {
	var req request.EmailLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ParamsErrorResp(c, err)
		return
	}
	common.WrapResp(c)(u.userService.LoginWithEmail(c.Request.Context(), &req))
}
