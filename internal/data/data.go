package data

import (
	"context"

	"entgo.io/ent/dialect"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/internal/data/ent"
	"github.com/liushuangls/go-server-template/internal/data/ent/migrate"
	"github.com/liushuangls/go-server-template/pkg/entutil"

	_ "github.com/go-sql-driver/mysql"
)

type Data struct {
	db    *ent.Client
	redis *redis.Client
}

func NewData(entClient *ent.Client, redisCli *redis.Client) (*Data, func(), error) {

	d := &Data{
		db:    entClient,
		redis: redisCli,
	}
	return d, func() {
		_ = d.db.Close()
	}, nil
}

func NewEntClient(conf *configs.Config) (*ent.Client, error) {
	var (
		drv dialect.Driver
		err error
	)

	drv, err = entutil.NewDriver(&entutil.Config{
		Dialect:     conf.DB.Dialect,
		DSN:         conf.DB.DSN,
		MaxIdle:     conf.DB.MaxIdle,
		MaxActive:   conf.DB.MaxActive,
		MaxLifetime: conf.DB.MaxLifetime,
	})
	if err != nil {
		return nil, err
	}
	if conf.IsDebugMode() {
		drv = entutil.Debug(drv)
	}

	client := ent.NewClient(ent.Driver(drv))
	if conf.DB.AutoMigrate {
		if err := client.Schema.Create(context.Background(), migrate.WithForeignKeys(false)); err != nil {
			return nil, err
		}
	}
	return client, nil
}

func NewRedisClient(conf *configs.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func NewRedisLocker(rdb *redis.Client) *redislock.Client {
	return redislock.New(rdb)
}

func NewRedisLimiter(rdb *redis.Client) *redis_rate.Limiter {
	return redis_rate.NewLimiter(rdb)
}
