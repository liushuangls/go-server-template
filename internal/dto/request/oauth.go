package request

import (
	"github.com/liushuangls/go-server-template/internal/data/ent/useroauth"
)

type OAuthPlatform struct {
	Platform string `json:"platform" form:"platform" binding:"required,oneof=google microsoft apple"`
}

func (o *OAuthPlatform) GetPlatform() useroauth.Platform {
	return useroauth.Platform(o.Platform)
}

type OAuthCodeURLReq struct {
	OAuthPlatform
}

type OAuthCallbackReq struct {
	OAuthPlatform
	Code  string `json:"code" form:"code" binding:"required,gte=1,lte=256"`
	State string `json:"state" form:"state" binding:"required,gte=1,lte=256"`
}
