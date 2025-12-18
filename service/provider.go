package service

import (
	"github.com/google/wire"

	"github.com/liushuangls/go-server-template/data"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Options), "*"),
	NewHealthService,
)

type Options struct {
	UserRepo *data.UserRepo
}
