package service

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"

	"github.com/liushuangls/go-server-template/internal/data"
	"github.com/liushuangls/go-server-template/internal/pkg"
	"github.com/liushuangls/go-server-template/internal/pkg/publicoss"
	"github.com/liushuangls/go-server-template/pkg/jwt"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Options), "*"),
	NewUserService,
)

type Options struct {
	Jwt          *jwt.JWT
	Redis        *redis.Client
	OAuthClients *pkg.OauthClients
	Avatar       *publicoss.Avatar

	UserRepo      *data.UserRepo
	UserOauthRepo *data.UserOAuthRepo
}
