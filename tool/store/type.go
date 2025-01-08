package store

import (
	"context"
	"time"
)

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	Del(ctx context.Context, key ...string) (bool, error)
}
