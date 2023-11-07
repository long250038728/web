package redis

import (
	"context"
	"testing"
	"time"
)

var addr = "43.139.51.99:32088"
var password = "zby123456"
var db = 0

var key = "z"
var value = "100"

func TestNewRedisCache(t *testing.T) {
	cache := NewRedisCache(addr, password, db)
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

	del, err := cache.Del(ctx, key)
	t.Log(del)
	t.Log(err)
}
