package limiter

import "context"

type Limiter interface {
	Allow(context context.Context, key string) (bool, error)
}
