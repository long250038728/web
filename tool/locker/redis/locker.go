package redis

import (
	"context"
	"errors"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/locker"
)

type Locker struct {
	client cache.Cache
}

var LockIsNotSelf = errors.New("lock is not self")

func NewRedisLocker(client cache.Cache) locker.Locker {
	return &Locker{
		client: client,
	}
}

func (l *Locker) Lock(context context.Context, key string) (bool, error) {
	return l.client.SetNX(context, key, key, 0)
}

func (l *Locker) UnLock(context context.Context, key string) (bool, error) {
	value, err := l.client.GetDel(context, key)
	if err != nil {
		return false, err
	}
	if key != value {
		return false, LockIsNotSelf
	}
	return true, nil
}

func (l *Locker) Del(context context.Context, key string) (bool, error) {
	return l.client.Del(context, key)
}
