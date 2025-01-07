package store

import (
	"context"
	"github.com/coocood/freecache"
	"github.com/long250038728/web/tool/persistence/redis"
	"time"
)

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	Del(ctx context.Context, key ...string) (bool, error)
}

func NewLocalStore(size int) Store {
	return &localStore{freecache.NewCache(size)}
}

func NewRedisStore(r redis.Redis) Store {
	return &redisStore{r}
}
