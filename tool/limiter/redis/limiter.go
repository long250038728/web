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

func (l *Limiter) Allow(ctx context.Context, key string) error {
	script := `
		if (redis.call("SETNX",KEYS[1],ARGV[1]) == 1) then
			redis.call("EXPIRE",KEYS[1],ARGV[2]);
		end
		return redis.call("incr",KEYS[1]);
	`
	data, err := l.client.Eval(ctx, script, []string{key}, 0, l.expiration)
	if err != nil {
		return err
	}
	num := data.(int64)
	if l.times >= num {
		return nil
	}
	return limiter.ErrorLimiter
}
