package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/liushuangls/go-server-template/internal/data"
	"github.com/liushuangls/go-server-template/internal/routes/common"
	"github.com/liushuangls/go-server-template/pkg/ecode"
	"github.com/liushuangls/go-server-template/pkg/jwt"
)

func TokenAuth(mustLogin bool, jwt *jwt.JWT, userRepo *data.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			if mustLogin {
				common.ErrorResp(c, ecode.InvalidToken)
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
			common.ErrorResp(c, ecode.InvalidToken)
			return
		}
		c.Set(common.CurrentUserInfoKey, u)
	}
}
