package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/liushuangls/go-server-template/internal/dto/request"
	"github.com/liushuangls/go-server-template/internal/routes/common"
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
		user.GET("/oauth/auth_code_url", u.getOAuthCodeUrl)
		user.GET("/oauth/callback", u.oauthCallback)
	}
}

func (u *UserRoute) getOAuthCodeUrl(c *gin.Context) {
	var req request.OAuthCodeURLReq
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ParamsErrorResp(c, err)
		return
	}
	common.WrapResp(c)(u.UserService.GetOAuthCodeURL(c.Request.Context(), &req))
}

func (u *UserRoute) oauthCallback(c *gin.Context) {
	var req request.OAuthCallbackReq
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ParamsErrorResp(c, err)
		return
	}
	common.WrapResp(c)(u.UserService.OAuthCallback(c.Request.Context(), common.MustGetIPInfo(c), &req))
}
