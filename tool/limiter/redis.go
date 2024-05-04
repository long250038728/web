package limiter

import (
	"context"
	"errors"
	"github.com/long250038728/web/tool/cache"
	"time"
)

type cacheLimiter struct {
	client     cache.Cache
	expiration time.Duration
	times      int64
}

func NewCacheLimiter(client cache.Cache, expiration time.Duration, times int64) Limiter {
	return &cacheLimiter{
		client:     client,
		expiration: expiration,
		times:      times,
	}
}

func (l *cacheLimiter) Allow(ctx context.Context, key string) error {
	script := `
		if (redis.call("SETNX",KEYS[1],ARGV[1]) == 1) then
			redis.call("PEXPIRE",KEYS[1],ARGV[2]);
		end
		return redis.call("incr",KEYS[1]);
	`
	if l.expiration < time.Millisecond {
		return errors.New("expiration time is error")
	}

	data, err := l.client.Eval(ctx, script, []string{key}, 0, l.expiration/time.Millisecond)
	if err != nil {
		return err
	}
	if l.times >= data.(int64) {
		return nil
	}
	return errors.New("API Limiter")
}
