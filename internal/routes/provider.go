package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/google/wire"
	"github.com/ip2location/ip2location-go/v9"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/internal/clients"
	v1 "github.com/liushuangls/go-server-template/internal/routes/v1"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Options), "*"),
	wire.Struct(new(v1.Options), "*"),
	NewEngine,
	NewHttpEngine,
	v1.NewUserRoute,
)

type Options struct {
	Router  *gin.Engine
	Conf    *configs.Config
	Limiter *redis_rate.Limiter
	IPDB    *ip2location.DB
	HashID  *clients.HashID

	User *v1.UserRoute
}
