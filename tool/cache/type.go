package cache

import (
	"context"
	"time"
)

type Cache interface {
	IsExists(ctx context.Context, key string) (bool, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	PTTL(ctx context.Context, key string) (time.Duration, error)
	Get(ctx context.Context, key string) (string, error)

	Set(ctx context.Context, key string, value string) (bool, error)
	SetEX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)

	Incr(ctx context.Context, key string) (int64, error)

	Del(ctx context.Context, key ...string) (bool, error)
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
}
