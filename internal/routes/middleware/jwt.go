package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/liushuangls/go-server-template/internal/data"
	entSchema "github.com/liushuangls/go-server-template/internal/data/ent"
	"github.com/liushuangls/go-server-template/internal/routes/common"
	"github.com/liushuangls/go-server-template/pkg/ecode"
	"github.com/liushuangls/go-server-template/pkg/jwt"
)

const currentUserInfo = "current-user-info"

func TokenAuth(mustLogin bool, jwt *jwt.JWT, userRepo *data.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			if mustLogin {
				common.ErrorResp(c, ecode.InvalidToken.WithCause(fmt.Errorf("no token")))
			}
			return
		}
		claims, err := jwt.ParseToken(token)
		if err != nil {
			if mustLogin {
				common.ErrorResp(c, err)
			}
			return
		}
		u, err := userRepo.FindByID(c.Request.Context(), claims.UserID)
		if err != nil {
			common.ErrorResp(c, ecode.InvalidToken.WithCause(fmt.Errorf("query user: %s", err)))
			return
		}
		c.Set(currentUserInfo, u)
	}
}

func GetCurrentUserInfo(c *gin.Context) *entSchema.User {
	val, exist := c.Get(currentUserInfo)
	if !exist {
		return nil
	}
	return val.(*entSchema.User)
}

func MustGetCurrentUserInfo(c *gin.Context) *entSchema.User {
	return c.MustGet(currentUserInfo).(*entSchema.User)
}
