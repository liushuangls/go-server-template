package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors(allowAll bool, allowOrigins ...string) gin.HandlerFunc {
	var config cors.Config
	if allowAll {
		config = cors.Config{
			AllowMethods:           []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:           []string{"*"},
			AllowCredentials:       false,
			MaxAge:                 12 * time.Hour,
			AllowWebSockets:        true,
			AllowWildcard:          true,
			AllowBrowserExtensions: true,
			AllowAllOrigins:        true,
		}
	} else {
		config = cors.DefaultConfig()
		config.AllowOrigins = allowOrigins
	}
	return cors.New(config)
}
