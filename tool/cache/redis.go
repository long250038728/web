package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var _ Cache = &Redis{}

type Config struct {
	Addr     string
	Password string
	Db       int
}

type Redis struct {
	client redis.Cmdable
}

func NewRedisCache(config *Config) Cache {
	////哨兵机制客户端init方法
	//redisClient := redis.NewFailoverClient(&redis.FailoverOptions{
	//	MasterName:"master",
	//	SentinelAddrs:[]string{"xxx1:6379","xxx2:6379","xxx3:6379"},
	//})
	redisClient := redis.NewClient(&redis.Options{Addr: config.Addr, Password: config.Password, DB: config.Db})
	return &Redis{
		client: redisClient,
	}
}

func (r *Redis) IsExists(ctx context.Context, key string) (bool, error) {
	res, err := r.client.Exists(ctx, key).Result()
	return res > 0, err
}

func (r *Redis) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

func (r *Redis) PTTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.PTTL(ctx, key).Result()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *Redis) Set(ctx context.Context, key string, value string) (bool, error) {
	result, err := r.client.Set(ctx, key, value, 0).Result()
	return result == "OK", err
}

func (r *Redis) SetEX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	result, err := r.client.SetEX(ctx, key, value, expiration).Result()
	return result == "OK", err
}

func (r *Redis) SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

func (r *Redis) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *Redis) Del(ctx context.Context, key ...string) (bool, error) {
	res, err := r.client.Del(ctx, key...).Result()
	return res > 0, err
}

func (r *Redis) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(ctx, script, keys, args...).Result()
}
