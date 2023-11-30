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

func NewRedisCache(addr, password string, db int) cache.Cache {
	redisClient := redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})
	return &Cache{
		client: redisClient,
	}
}

func (r *Cache) IsExists(context context.Context, key string) (bool, error) {
	res, err := r.client.Exists(context, key).Result()
	return res > 0, err
}

func (r *Cache) TTL(context context.Context, key string) (time.Duration, error) {
	return r.client.TTL(context, key).Result()
}

func (r *Cache) Get(context context.Context, key string) (string, error) {
	return r.client.Get(context, key).Result()
}

func (r *Cache) Set(context context.Context, key string, value string) (bool, error) {
	result, err := r.client.Set(context, key, value, 0).Result()
	return result == "OK", err
}

func (r *Cache) SetEX(context context.Context, key string, value string, expiration time.Duration) (bool, error) {
	result, err := r.client.SetEX(context, key, value, expiration).Result()
	return result == "OK", err
}

func (r *Cache) SetNX(context context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return r.client.SetNX(context, key, value, expiration).Result()
}

func (r *Cache) Incr(context context.Context, key string) (int64, error) {
	return r.client.Incr(context, key).Result()
}

func (r *Cache) Del(context context.Context, key ...string) (bool, error) {
	res, err := r.client.Del(context, key...).Result()
	return res > 0, err
}
