package global

import (
	"context"

	"github.com/gin-contrib/sessions"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	RSS sessions.Store
)

func SetupRedis() {
	opts, err := redis.ParseURL(CFG.RedisURI)
	if err != nil {
		panic(err)
	}

	RDB = redis.NewClient(opts)
	if err := redisotel.InstrumentTracing(RDB); err != nil {
		panic(err)
	}
	store, err := redisstore.NewRedisStore(context.Background(), RDB)
	if err != nil {
		panic(err)
	}
	RSS = SessionsStore{store}
}

type SessionsStore struct {
	*redisstore.RedisStore
}

func (s SessionsStore) Options(options sessions.Options) {
	s.RedisStore.Options(*options.ToGorillaOptions())
}
