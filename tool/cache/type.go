package cache

import (
	"context"
	"time"
)

type Cache interface {
	Exists(context context.Context, key string) (bool, error)
	TTL(context context.Context, key string) (time.Duration, error)
	Get(context context.Context, key string) (string, error)

	Set(context context.Context, key string, value string) (bool, error)
	SetEX(context context.Context, key string, value string, expiration time.Duration) (bool, error)
	SetNX(context context.Context, key string, value string, expiration time.Duration) (bool, error)
}
