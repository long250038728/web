package limiter

import (
	"context"
	"github.com/long250038728/web/tool/persistence/redis"
	"strconv"
)

var _ Limiter = &redisLimiter{}

type redisLimiter struct {
	client redis.Redis
}

func (l *redisLimiter) Get(ctx context.Context, key string) (int64, error) {
	cnt, err := l.client.Get(ctx, key)
	if err != nil {
		return 0, nil
	}
	return strconv.ParseInt(cnt, 10, 0)
}

func (l *redisLimiter) Incr(ctx context.Context, key string) (int64, error) {
	return l.client.Incr(ctx, key)
}

func (l *redisLimiter) Decr(ctx context.Context, key string) (int64, error) {
	return l.client.Decr(ctx, key)
}
