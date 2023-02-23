package common

import (
	"github.com/gin-gonic/gin"
)

func GetRealIP(c *gin.Context) string {
	ip := c.GetHeader("X-Real-IP")
	if ip == "" {
		ip = c.GetHeader("X-Forwarded-For")
	}
	if ip == "" {
		ip = c.ClientIP()
	}
	return ip
}
