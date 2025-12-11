package cmd

import (
	"github.com/google/wire"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/internal/crontab"
	"github.com/liushuangls/go-server-template/internal/routes"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Options), "*"),
	NewApp,
)

type Options struct {
	Config *configs.Config
	Http   *routes.HttpEngine
	Cron   *crontab.Client
}
