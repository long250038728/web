package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var _ Cache = &redisClient{}

type Config struct {
	Address  string `json:"address" yaml:"address"`
	Password string `json:"password" yaml:"password"`
	Db       int    `json:"db" yaml:"db"`
}

type redisClient struct {
	client *redis.Client
}

func NewRedis(config *Config) Cache {
	////哨兵机制客户端init方法
	//redisClient := cache.NewFailoverClient(&cache.FailoverOptions{
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

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	if expiration == 0 {
		result, err := r.client.Set(ctx, key, value, 0).Result()
		return result == "OK", err
	}
	result, err := r.client.SetEX(ctx, key, value, expiration).Result()
	return result == "OK", err
}

func (r *redisClient) Del(ctx context.Context, key ...string) (bool, error) {
	res, err := r.client.Del(ctx, key...).Result()
	return res > 0, err
}

func (r *redisClient) LPush(ctx context.Context, key string, value string) (bool, error) {
	res, err := r.client.LPush(ctx, key, value).Result()
	return res > 0, err
}

func (r *redisClient) SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

func (r *redisClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(ctx, script, keys, args...).Result()
}

func (r *redisClient) Publish(ctx context.Context, channel, message string) (bool, error) {
	res, err := r.client.Publish(ctx, channel, message).Result()
	return res > 0, err
}

func (r *redisClient) Subscribe(ctx context.Context, channel string, f func(message string)) {
	sub := r.client.Subscribe(ctx, channel)
	for msg := range sub.Channel() {
		f(msg.Payload)
	}
}
