package redis

import (
	"context"
	"fmt"
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

func TestNewRedisCache(t *testing.T) {
	cache := NewRedisCache(config)
	ctx := context.Background()

	ok, err := cache.Set(ctx, key, value)
	t.Log(ok)
	t.Log(err)

	val, err := cache.Get(ctx, key)
	t.Log(val)
	t.Log(err)

	ttl, err := cache.TTL(ctx, key)
	t.Log(ttl)
	t.Log(err)

	valueEx, err := cache.SetEX(ctx, key, value, time.Second)
	t.Log(valueEx)
	t.Log(err)

	valueNx, err := cache.SetNX(ctx, key, value, 0)
	t.Log(valueNx)
	t.Log(err)

	v, err := cache.Get(ctx, key)
	t.Log(v)
	t.Log(err)

	del, err := cache.Del(ctx, key)
	t.Log(del)
	t.Log(err)
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

	// 使用客户端进行操作
	ctx := context.Background()

	// 示例：设置键值对
	err := client.Set(ctx, "example_key", "example_value", 0).Err()
	if err != nil {
		fmt.Println("Error setting value:", err)
		return
	}
}
