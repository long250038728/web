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

func (l *Locker) Lock(ctx context.Context, key, identification string) (bool, error) {
	return l.client.SetNX(ctx, key, identification, 0)
}

func (l *Locker) UnLock(ctx context.Context, key, identification string) (bool, error) {
	script := `
		if (redis.call("get",KEYS[1]) == ARGV[1]) then
			return redis.call("del",KEYS[1]);
		else
			return 0;
		end
	`
	data, err := l.client.Eval(ctx, script, []string{key}, identification)
	if err != nil {
		return false, err
	}
	if data.(int64) == 0 {
		return false, locker.ErrorIdentification
	}
	return true, nil
}
