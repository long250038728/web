package redis

import (
	"context"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/locker"
)

type Locker struct {
	client cache.Cache
}

func NewRedisLocker(client cache.Cache) locker.Locker {
	return &Locker{
		client: client,
	}
}

func (l *Locker) Lock(context context.Context, key string) (bool, error) {
	return l.client.SetNX(context, key, key, 0)
}

func (l *Locker) UnLock(context context.Context, key string) (bool, error) {
	return l.client.Del(context, key)
}
