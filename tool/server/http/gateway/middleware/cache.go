package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/app_error"
	"github.com/long250038728/web/tool/authorization"
	"github.com/long250038728/web/tool/persistence/cache"
	"github.com/long250038728/web/tool/server/http/gateway"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"reflect"
	"time"
)

type middlewareInfo struct {
	claims     bool
	expiration time.Duration
}

func NewDefaultMiddlewareInfo() *middlewareInfo {
	return &middlewareInfo{expiration: time.Minute}
}

type Opt func(setting *middlewareInfo)

func Expiration(expiration time.Duration) Opt {
	return func(setting *middlewareInfo) {
		setting.expiration = expiration
	}
}
func Claims(claims bool) Opt {
	return func(setting *middlewareInfo) {
		setting.claims = claims
	}
}

// =======================================================================================================
func getKey(c *gin.Context, isClaims bool, keys []string, requestInfo map[string]any) (string, bool) {
	key := c.Request.URL.Path
	dataIsAllMatch := true

	if isClaims {
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

// =======================================================================================================

// Locker 解决并发锁
func Locker(c *gin.Context, client cache.Cache, keys []string, opts ...Opt) gateway.ServerInterceptor {
	return func(ctx context.Context, requestInfo map[string]any, request any, handler gateway.Handler) (resp any, err error) {
		setting := NewDefaultMiddlewareInfo()
		for _, opt := range opts {
			opt(setting)
		}
		expiration := setting.expiration
		key, dataIsAllMatch := getKey(c, setting.claims, keys, requestInfo)
		if !dataIsAllMatch {
			return nil, app_error.Vaildate
		}
		ok, err := client.SetNX(ctx, key, "1", expiration)
		if err != nil || !ok {
			return nil, app_error.ApiTooManyRequests
		}
		defer func() {
			_, _ = client.Del(ctx, key)
		}()
		return handler(ctx, request)
	}
}

// Cache  接口缓存
func Cache(c *gin.Context, client cache.Cache, keys []string, opts ...Opt) gateway.ServerInterceptor {
	return func(ctx context.Context, requestInfo map[string]any, request any, handler gateway.Handler) (resp any, err error) {
		setting := NewDefaultMiddlewareInfo()
		for _, opt := range opts {
			opt(setting)
		}
		expiration := setting.expiration
		key, dataIsAllMatch := getKey(c, setting.claims, keys, requestInfo)

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
				_, _ = client.Set(ctx, key, string(b), expiration)
			}
			return cacheData, err
		}

		return resp, err
	}
}
