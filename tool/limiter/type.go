package limiter

import (
	"context"
)

type Limiter interface {
	Get(ctx context.Context, key string) (int64, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	Allow(ctx context.Context, key string, num int64) (bool, error)
}
