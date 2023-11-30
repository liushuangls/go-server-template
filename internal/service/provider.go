package service

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"

	"github.com/liushuangls/go-server-template/internal/clients"
	"github.com/liushuangls/go-server-template/internal/clients/publicoss"
	"github.com/liushuangls/go-server-template/internal/data"
	"github.com/liushuangls/go-server-template/pkg/jwt"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Options), "*"),
	NewUserService,
)

type Options struct {
	Jwt          *jwt.JWT
	Redis        *redis.Client
	OAuthClients *clients.OauthClients
	Avatar       *publicoss.Avatar
	HashID       *clients.HashID

	UserRepo      *data.UserRepo
	UserOauthRepo *data.UserOAuthRepo
}
