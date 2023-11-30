package redis

import (
	"context"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/limiter"
	"time"
)

type Limiter struct {
	client     cache.Cache
	expiration time.Duration
	times      int64
}

func NewRedisLimiter(client cache.Cache, expiration time.Duration, times int64) limiter.Limiter {
	return &Limiter{
		client:     client,
		expiration: expiration,
		times:      times,
	}
}

func (l *Limiter) Allow(ctx context.Context, key string) (bool, error) {
	//如果不存在就插入一个数据
	_, err := l.client.SetNX(ctx, key, "0", l.expiration)
	if err != nil {
		return false, err
	}

	//++1
	num, err := l.client.Incr(ctx, key)
	if err != nil {
		return false, err
	}

	//配置值大于等于自增数量，允许执行
	return l.times >= num-1, nil
}
