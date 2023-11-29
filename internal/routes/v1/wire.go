package v1

import (
	"github.com/bsm/redislock"
	"github.com/go-redis/redis_rate/v10"

	"github.com/liushuangls/go-server-template/internal/data"
	"github.com/liushuangls/go-server-template/internal/service"
	"github.com/liushuangls/go-server-template/pkg/jwt"
)

type Options struct {
	Limiter *redis_rate.Limiter
	Locker  *redislock.Client
	Jwt     *jwt.JWT

	UserRepo *data.UserRepo

	UserService *service.UserService
}
