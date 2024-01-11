package limiter

import (
	"context"
	"errors"
)

var ErrorLimiter = errors.New("api limiter")

type Limiter interface {
	Allow(ctx context.Context, key string) error
}
