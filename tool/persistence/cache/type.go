package cache

import (
	"context"
	"time"
)

type Cache interface {
	IsExists(ctx context.Context, key string) (bool, error)
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	Del(ctx context.Context, key ...string) (bool, error)
	MQ
	Locker
}

// MQ 用于消息队列
type MQ interface {
	LPush(ctx context.Context, key string, value string) (bool, error)

	Publish(ctx context.Context, channel, message string) (bool, error)
	Subscribe(ctx context.Context, channel string, f func(message string))
}

// Locker 锁
type Locker interface {
	SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
}
