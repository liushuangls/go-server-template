package controller

import (
	"github.com/bsm/redislock"
	"github.com/go-redis/redis_rate/v10"

	"github.com/liushuangls/go-server-template/service"
)

type Options struct {
	Limiter *redis_rate.Limiter
	Locker  *redislock.Client

	HealthService *service.HealthService
}
