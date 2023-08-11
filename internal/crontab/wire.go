package crontab

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Options), "*"),
	NewClient,
)

type Options struct {
}
