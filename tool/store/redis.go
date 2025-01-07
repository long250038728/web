package store

import (
	"context"
	"github.com/long250038728/web/tool/persistence/redis"
	"time"
)

type redisStore struct {
	r redis.Redis
}

func (l *redisStore) Get(ctx context.Context, key string) (string, error) {
	return l.r.Get(ctx, key)
}

func (l *redisStore) Set(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return l.r.SetEX(ctx, key, value, expiration)
}

func (l *redisStore) Del(ctx context.Context, key ...string) (bool, error) {
	return l.r.Del(ctx, key...)
}
