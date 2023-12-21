package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

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
	needAuth := middleware.TokenAuth(true, u.Jwt, u.UserRepo)
	{
		user.GET("/oauth/auth_code_url", u.getOAuthCodeUrl)
		user.GET("/oauth/callback", u.oauthCallback)

		user.POST("/oauth/oauth_one_tap", u.oauthOneTap)

		user.GET("/info", needAuth, u.userInfo)
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

func (u *UserRoute) userInfo(c *gin.Context) {
	user := common.MustGetCurrentUserInfo(c)
	common.WrapResp(c)(u.UserService.UserInfo(c.Request.Context(), user))
}

func (u *UserRoute) oauthOneTap(c *gin.Context) {
	var req request.OauthOneTapReq
	req.Platform = "google"
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		common.ParamsErrorResp(c, err)
		return
	}

	common.WrapResp(c)(u.UserService.OAuthOneTap(c.Request.Context(), common.MustGetIPInfo(c), &req))
}
