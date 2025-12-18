package cmd

import (
	"github.com/google/wire"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/crontab"
	"github.com/liushuangls/go-server-template/routes"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Options), "*"),
	NewApp,
	NewDefaultSlog,
)

type Options struct {
	Config *configs.Config
	Http   *routes.HttpEngine
	Cron   *crontab.Client
}
