package locker

import (
	"context"
	"github.com/long250038728/web/tool/cache"
	"time"
)

type redis struct {
	client cache.Cache
	key,
	identification string
	time time.Duration
}

func NewRedisLocker(client cache.Cache, key, identification string, RefreshTime time.Duration) Locker {
	return &redis{
		client:         client,
		key:            key,
		identification: identification,
		time:           RefreshTime,
	}
}

func (l *redis) Lock(ctx context.Context) error {
	ok, err := l.client.SetNX(ctx, l.key, l.identification, l.time)
	if err != nil {
		return err
	}
	if !ok {
		return ErrorIsExists
	}
	return nil
}

func (l *redis) UnLock(ctx context.Context) error {
	script := `
		if (redis.call("get",KEYS[1]) == ARGV[1]) then
			return redis.call("del",KEYS[1]);
		else
			return 0;
		end
	`
	data, err := l.client.Eval(ctx, script, []string{l.key}, l.identification)
	if err != nil {
		return err
	}
	if data.(int64) == 0 {
		return ErrorIdentification
	}
	return nil
}

func (l *redis) Refresh(ctx context.Context) error {
	script := `
		if (redis.call("get",KEYS[1]) == ARGV[1]) then
			return redis.call("expire",KEYS[1],ARGV[2]);
		else
			return 0;
		end
	`
	data, err := l.client.Eval(ctx, script, []string{l.key}, l.identification, l.time)
	if err != nil {
		return err
	}
	if data.(int64) == 0 {
		return ErrorIdentification
	}
	return nil
}

func (l *redis) AutoRefresh(ctx context.Context) error {
	retryTimes := 3
	retry := 0

	t := l.time - time.Microsecond*500 //往前推个半秒，避免时间到了此时因网络延迟，续约已经找不到了
	//续约时间小于等于0则无需续约
	if t <= 0 {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			// 超时控制
			return ctx.Err()

		case <-time.After(t):
			//重试
			for {
				err := l.Refresh(ctx)
				//成功续约退出重试循环
				if err == nil {
					retry = 0
					break
				}

				//如果续约失败，那续约直到到达重试次数
				if retry >= retryTimes {
					retry = 0
					return ErrorOverRetry
				}
				retry++
			}
		}
	}
}

func (l *redis) Close() error {
	return nil
}
