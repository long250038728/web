package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var _ Redis = &redisClient{}

type Config struct {
	Address  string `json:"address" yaml:"address"`
	Password string `json:"password" yaml:"password"`
	Db       int    `json:"db" yaml:"db"`
}

type redisClient struct {
	client redis.Cmdable
}

func NewRedis(config *Config) Redis {
	////哨兵机制客户端init方法
	//redisClient := redis.NewFailoverClient(&redis.FailoverOptions{
	//	MasterName:"master",
	//	SentinelAddrs:[]string{"xxx1:6379","xxx2:6379","xxx3:6379"},
	//})
	return &redisClient{
		client: redis.NewClient(&redis.Options{Addr: config.Address, Password: config.Password, DB: config.Db}),
	}
}

func (r *redisClient) IsExists(ctx context.Context, key string) (bool, error) {
	res, err := r.client.Exists(ctx, key).Result()
	return res > 0, err
}

func (r *redisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

func (r *redisClient) PTTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.PTTL(ctx, key).Result()
}

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisClient) Set(ctx context.Context, key string, value string) (bool, error) {
	result, err := r.client.Set(ctx, key, value, 0).Result()
	return result == "OK", err
}

func (r *redisClient) SetEX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	result, err := r.client.SetEX(ctx, key, value, expiration).Result()
	return result == "OK", err
}

func (r *redisClient) SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

func (r *redisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *redisClient) Decr(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

func (r *redisClient) Del(ctx context.Context, key ...string) (bool, error) {
	res, err := r.client.Del(ctx, key...).Result()
	return res > 0, err
}

func (r *redisClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(ctx, script, keys, args...).Result()
}
