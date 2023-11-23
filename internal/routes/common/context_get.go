package common

import (
	"github.com/gin-gonic/gin"

	"github.com/liushuangls/go-server-template/internal/data/ent"
	"github.com/liushuangls/go-server-template/internal/dto/request"
)

const (
	IpInfoKey          = "client-ip-info"
	CurrentUserInfoKey = "current-user-info"
)

func MustGetIPInfo(c *gin.Context) *request.IPInfo {
	return c.MustGet(IpInfoKey).(*request.IPInfo)
}

func GetIPInfo(c *gin.Context) *request.IPInfo {
	val, exist := c.Get(IpInfoKey)
	if !exist {
		return nil
	}
	return val.(*request.IPInfo)
}

func GetCurrentUserInfo(c *gin.Context) *ent.User {
	val, exist := c.Get(CurrentUserInfoKey)
	if !exist {
		return nil
	}
	return val.(*ent.User)
}

func MustGetCurrentUserInfo(c *gin.Context) *ent.User {
	return c.MustGet(CurrentUserInfoKey).(*ent.User)
}
