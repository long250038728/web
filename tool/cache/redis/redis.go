package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/long250038728/web/tool/cache"
	"time"
)

var _ cache.Cache = &Cache{}

type Cache struct {
	client redis.Cmdable
}

func NewRedisCache(config *Config) cache.Cache {
	////哨兵机制客户端init方法
	//redisClient := redis.NewFailoverClient(&redis.FailoverOptions{
	//	MasterName:"master",
	//	SentinelAddrs:[]string{"xxx1:6379","xxx2:6379","xxx3:6379"},
	//})
	redisClient := redis.NewClient(&redis.Options{Addr: config.Addr, Password: config.Password, DB: config.Db})
	return &Cache{
		client: redisClient,
	}
}

func (r *Cache) IsExists(ctx context.Context, key string) (bool, error) {
	res, err := r.client.Exists(ctx, key).Result()
	return res > 0, err
}

func (r *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

func (r *Cache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *Cache) Set(ctx context.Context, key string, value string) (bool, error) {
	result, err := r.client.Set(ctx, key, value, 0).Result()
	return result == "OK", err
}

func (r *Cache) SetEX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	result, err := r.client.SetEX(ctx, key, value, expiration).Result()
	return result == "OK", err
}

func (r *Cache) SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

func (r *Cache) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *Cache) Del(ctx context.Context, key ...string) (bool, error) {
	res, err := r.client.Del(ctx, key...).Result()
	return res > 0, err
}

func (r *Cache) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(ctx, script, keys, args...).Result()
}
