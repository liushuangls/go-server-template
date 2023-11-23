package response

import (
	"github.com/liushuangls/go-server-template/pkg/jwt"
)

type UserInfo struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	NickName   string `json:"nickname"`
	RegisterAt int64  `json:"register_at"`
}

type UserLoginInfo struct {
	UserInfo
	AccessToken *jwt.Token `json:"access_token"`
}
