package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/authorization"
	"github.com/long250038728/web/tool/persistence/redis"
	"github.com/long250038728/web/tool/server/http/gateway"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"reflect"
	"time"
)

type cacheSetting struct {
	isClaims   bool
	expiration time.Duration
}

func NewDefaultCacheSetting() *cacheSetting {
	return &cacheSetting{
		expiration: time.Minute,
	}
}

type CacheOpt func(setting *cacheSetting)

func SetExpiration(expiration time.Duration) CacheOpt {
	return func(setting *cacheSetting) {
		setting.expiration = expiration
	}
}
func SetIsClaims(isClaims bool) CacheOpt {
	return func(setting *cacheSetting) {
		setting.isClaims = isClaims
	}
}

//=======================================================================================================

func getCacheKey(c *gin.Context, setting *cacheSetting, keys []string, requestInfo map[string]any) (string, bool) {
	key := c.Request.URL.Path
	dataIsAllMatch := true

	if setting.isClaims {
		c, err := authorization.GetClaims(c.Request.Context())
		if err != nil {
			dataIsAllMatch = false
			return key, dataIsAllMatch
		}
		key += ":" + fmt.Sprintf("%v", c.Id)
	}

	for _, k := range keys {
		// 参数
		d, ok := requestInfo[k]
		if !ok || reflect.DeepEqual(d, reflect.Zero(reflect.TypeOf(d)).Interface()) {
			dataIsAllMatch = false
			return key, dataIsAllMatch
		}
		key += ":" + fmt.Sprintf("%v", d)
	}

	return key, dataIsAllMatch
}

func Cache(c *gin.Context, client redis.Redis, keys []string, opts ...CacheOpt) gateway.ServerInterceptor {
	return func(ctx context.Context, requestInfo map[string]any, request any, handler gateway.Handler) (resp any, err error) {
		setting := NewDefaultCacheSetting()
		for _, opt := range opts {
			opt(setting)
		}
		expiration := setting.expiration
		key, dataIsAllMatch := getCacheKey(c, setting, keys, requestInfo)

		if !dataIsAllMatch {
			return handler(ctx, request)
		}

		// 读取缓存
		// 缓存读到数据 && 缓存转换为 Response对象 没有错误
		// 直接返回 Response对象
		b, err := client.Get(ctx, key)
		if err == nil && len(b) > 0 {
			var cacheData *gateway.Response
			if err := json.Unmarshal([]byte(b), &cacheData); err == nil {

				span := opentelemetry.NewSpan(ctx, "middle_cache: "+key)
				defer span.Close()
				span.SetAttribute("middle_cache", true)

				return cacheData, nil
			}
		}

		resp, err = handler(ctx, request)

		// 写入缓存
		// 转换成 Response对象
		if err == nil {
			cacheData := gateway.NewResponse(resp, err)
			if b, err := json.Marshal(cacheData); err == nil {
				_, _ = client.SetEX(ctx, key, string(b), expiration)
			}
			return cacheData, err
		}

		return resp, err
	}
}
