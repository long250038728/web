package store

import (
	"context"
	"github.com/long250038728/web/tool/persistence/cache"
	"time"
)

type redisStore struct {
	r cache.Cache
}

func NewRedisStore(r cache.Cache) Store {
	return &redisStore{r}
}

func (s *redisStore) Get(ctx context.Context, key string) (string, error) {
	return s.r.Get(ctx, key)
}

func (s *redisStore) Set(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return s.r.Set(ctx, key, value, expiration)
}

func (s *redisStore) Del(ctx context.Context, key ...string) (bool, error) {
	return s.r.Del(ctx, key...)
}

func (s *redisStore) Close() {

}
