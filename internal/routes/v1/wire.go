package v1

import (
	"github.com/bsm/redislock"
	"github.com/go-redis/redis_rate/v10"
	"go.uber.org/zap"

	"github.com/liushuangls/go-server-template/internal/service"
)

type Options struct {
	Log     *zap.SugaredLogger
	Limiter *redis_rate.Limiter
	Locker  *redislock.Client

	UserService *service.UserService
}
