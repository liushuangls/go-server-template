package service

import (
	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/liushuangls/go-server-template/internal/data"
	"github.com/liushuangls/go-server-template/pkg/jwt"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Options), "*"),
	NewUserService,
)

type Options struct {
	Log      *zap.SugaredLogger
	Jwt      *jwt.JWT
	UserRepo *data.UserRepo
}
