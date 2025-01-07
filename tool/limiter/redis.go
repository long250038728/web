package limiter

import (
	"context"
	"github.com/long250038728/web/tool/persistence/redis"
	"strconv"
)

var _ Limiter = &redisLimiter{}

type redisLimiter struct {
	client redis.Redis
}

func (l *redisLimiter) Get(ctx context.Context, key string) (int64, error) {
	cnt, err := l.client.Get(ctx, key)
	if err != nil {
		return 0, nil
	}
	return strconv.ParseInt(cnt, 10, 0)
}

func (l *redisLimiter) Incr(ctx context.Context, key string) (int64, error) {
	return l.client.Incr(ctx, key)
}

func (l *redisLimiter) Decr(ctx context.Context, key string) (int64, error) {
	return l.client.Decr(ctx, key)
}

func (l *redisLimiter) Allow(ctx context.Context, key string, num int64) (bool, error) {
	script := `
		local current = redis.call("get", KEYS[1])
		if not current then
			-- 如果键不存在，初始化为 1，代表允许
			current = 1
		else
			-- 确保 current 是数字类型
			current = tonumber(current)
		end
		
		if current < tonumber(ARGV[1]) then
			-- 如果当前值小于 num，增加计数。代表允许
			return redis.call("incr", KEYS[1]) 
		else
			-- 否则拒绝
			return 0
		end
	`
	data, err := l.client.Eval(ctx, script, []string{key}, num)
	return data.(int64) > 0, err
}
