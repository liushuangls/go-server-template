package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors(allowAll bool, allowOrigins ...string) gin.HandlerFunc {
	config := cors.DefaultConfig()
	if allowAll {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = allowOrigins
	}
	return cors.New(config)
}
