package crontab

import (
	"github.com/google/wire"
	"go.uber.org/zap"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Options), "*"),
	NewClient,
)

type Options struct {
	Log *zap.SugaredLogger
}
