package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"

	"github.com/liushuangls/go-server-template/internal/dto/request"
	"github.com/liushuangls/go-server-template/internal/routes/common"
	"github.com/liushuangls/go-server-template/internal/routes/middleware"
)

type UserRoute struct {
	Options
}

func NewUserRoute(opt Options) *UserRoute {
	return &UserRoute{opt}
}

func (u *UserRoute) RegisterRoute(router *gin.RouterGroup) {
	user := router.Group("/v1/user")
	{
		user.POST(
			"/login",
			middleware.RateLimitWithIP(u.Limiter, redis_rate.PerMinute(5), "login"),
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
	common.WrapResp(c)(u.UserService.LoginWithEmail(c.Request.Context(), &req))
}
