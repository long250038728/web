package limiter

import (
	"context"
	"github.com/long250038728/web/tool/persistence/redis"
)

type Limiter interface {
	Get(ctx context.Context, key string) (int64, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
}

func NewLocalLimiter() Limiter {
	return &localLimiter{data: map[string]int64{}}
}

func NewRedisLimiter(client redis.Redis) Limiter {
	return &redisLimiter{client: client}
}
