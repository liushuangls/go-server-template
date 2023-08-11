package service

import (
	"github.com/google/wire"

	"github.com/liushuangls/go-server-template/internal/data"
	"github.com/liushuangls/go-server-template/pkg/jwt"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Options), "*"),
	NewUserService,
)

type Options struct {
	Jwt      *jwt.JWT
	UserRepo *data.UserRepo
}
