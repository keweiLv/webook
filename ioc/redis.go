package ioc

import (
	"github.com/keweiLv/webook/env"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: env.Config.Redis.Addr,
	})
	return redisClient
}
