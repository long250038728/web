package limiter

import (
	"context"
)

type Store interface {
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
}

type Limiter interface {
	Allow(ctx context.Context, key string) error
}
