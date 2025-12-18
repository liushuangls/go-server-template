package middleware

import (
	"fmt"

	"github.com/go-redis/redis_rate/v10"
	"github.com/labstack/echo/v4"

	"github.com/liushuangls/go-server-template/pkg/ecode"
)

func RateLimitWithIP(limiter *redis_rate.Limiter, limit redis_rate.Limit, prefix string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := fmt.Sprintf("%s:%s", prefix, c.RealIP())
			result, err := limiter.Allow(c.Request().Context(), key, limit)
			if err != nil || result.Allowed == 0 {
				return ecode.TooManyRequest
			}
			return next(c)
		}
	}
}
