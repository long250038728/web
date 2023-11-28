package cache

import (
	"context"
	"time"
)

type Cache interface {
	IsExists(context context.Context, key string) (bool, error)
	TTL(context context.Context, key string) (time.Duration, error)
	Get(context context.Context, key string) (string, error)

	Set(context context.Context, key string, value string) (bool, error)
	SetEX(context context.Context, key string, value string, expiration time.Duration) (bool, error)
	SetNX(context context.Context, key string, value string, expiration time.Duration) (bool, error)

	Incr(context context.Context, key string) (int64, error)

	GetDel(context context.Context, key string) (string, error)
	Del(context context.Context, key ...string) (bool, error)
}
