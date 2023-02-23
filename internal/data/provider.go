package data

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewData, NewEntClient, NewRedisClient, NewRedisLocker, NewRedisLimiter, NewUserRepo)
