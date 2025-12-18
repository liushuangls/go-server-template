package routes

import (
	"github.com/go-redis/redis_rate/v10"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/routes/controller"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Options), "*"),
	wire.Struct(new(controller.Options), "*"),
	NewEcho,
	NewHttpEngine,
	controller.NewHealthRoute,
)

type Options struct {
	Router  *echo.Echo
	Conf    *configs.Config
	Limiter *redis_rate.Limiter

	Health *controller.HealthRoute
}
