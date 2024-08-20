package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/google/wire"
	"github.com/oschwald/maxminddb-golang"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/internal/clients"
	"github.com/liushuangls/go-server-template/internal/data"
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
	IPDB    *maxminddb.Reader
	HashID  *clients.HashID

	ServerLogRepo *data.ServerLogRepo

	User *v1.UserRoute
}
