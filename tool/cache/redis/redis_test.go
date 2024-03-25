package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

var config = &Config{
	Addr:     "43.139.51.99:32088",
	Password: "zby123456",
	Db:       0,
}

var key = "z"
var value = "100"

var cache = NewRedisCache(config)
var ctx = context.Background()

func TestNewRedisCache(t *testing.T) {
	t.Log(cache.Set(ctx, key, value))
	t.Log(cache.Get(ctx, key))
	t.Log(cache.TTL(ctx, key))
	t.Log(cache.PTTL(ctx, key))
	t.Log(cache.SetEX(ctx, key, value, time.Second))
	t.Log(cache.SetNX(ctx, key, value, 0))
	t.Log(cache.Del(ctx, key))
}
func TestSentinel(t *testing.T) {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: []string{"localhost:26379", "localhost:26380", "localhost:26381"},
		DB:            0,
		DialTimeout:   5 * time.Second,
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  5 * time.Second,
	})
	t.Log(client.Set(ctx, "example_key", "example_value", 0).Err())
}
