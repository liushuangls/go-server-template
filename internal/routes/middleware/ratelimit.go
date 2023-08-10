package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"

	"github.com/liushuangls/go-server-template/internal/routes/common"
	"github.com/liushuangls/go-server-template/pkg/ecode"
)

func RateLimitWithIP(limiter *redis_rate.Limiter, limit redis_rate.Limit, prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("%s:%s", prefix, common.GetRealIP(c))
		result, err := limiter.Allow(c.Request.Context(), key, limit)
		if err != nil || result.Allowed == 0 {
			common.ErrorResp(c, ecode.TooManyRequest)
			return
		}
	}
}
