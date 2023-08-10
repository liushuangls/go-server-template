package cmd

import (
	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/liushuangls/go-server-template/internal/crontab"
	"github.com/liushuangls/go-server-template/internal/routes"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Options), "*"),
	NewLogger,
	NewJwt,
	NewApp,
)

type Options struct {
	Log  *zap.SugaredLogger
	Http *routes.HttpEngine
	Cron *crontab.Client
}
