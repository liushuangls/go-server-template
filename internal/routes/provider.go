package routes

import (
	"github.com/go-redis/redis_rate/v10"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"

	"github.com/liushuangls/go-server-template/configs"
	v1 "github.com/liushuangls/go-server-template/internal/routes/controller"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Options), "*"),
	wire.Struct(new(v1.Options), "*"),
	NewEcho,
	NewHttpEngine,
	v1.NewHealthRoute,
)

type Options struct {
	Router  *echo.Echo
	Conf    *configs.Config
	Limiter *redis_rate.Limiter

	Health *v1.HealthRoute
}
