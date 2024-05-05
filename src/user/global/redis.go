package global

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func SetupRedis() {
	opts, err := redis.ParseURL(CFG.RedisURI)
	if err != nil {
		panic(err)
	}

	RDB = redis.NewClient(opts)
	if err = RDB.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}
