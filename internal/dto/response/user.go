package response

import (
	"github.com/liushuangls/go-server-template/pkg/jwt"
)

type UserInfo struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type UserLoginInfo struct {
	UserInfo
	AccessToken *jwt.Token `json:"access_token"`
}
