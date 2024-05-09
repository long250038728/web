package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/long250038728/web/tool/configurator"
	"testing"
	"time"
)

var key = "z"
var value = "100"
var cacheClient Cache

func init() {
	var conf Config
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/redis.yaml", &conf)
	cacheClient = NewRedisCache(&conf)
}

var ctx = context.Background()

func TestNewRedisCache(t *testing.T) {
	t.Log(cacheClient.Set(ctx, key, value))
	t.Log(cacheClient.Get(ctx, key))
	t.Log(cacheClient.TTL(ctx, key))
	t.Log(cacheClient.PTTL(ctx, key))
	t.Log(cacheClient.SetEX(ctx, key, value, time.Second))
	t.Log(cacheClient.SetNX(ctx, key, value, 0))
	t.Log(cacheClient.Del(ctx, key))
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
