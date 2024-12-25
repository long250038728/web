package gateway

import (
	"context"
	"fmt"
)

// Cache 示例中间件：缓存拦截器
func Cache() ServerInterceptor {
	return func(ctx context.Context, requestInfo map[string]any, request any, handler Handler) (resp any, err error) {
		fmt.Println("cache")
		// 检查缓存逻辑（省略实际实现）
		return handler(ctx, request)
	}
}

// Limit 示例中间件：限流拦截器
func Limit() ServerInterceptor {
	return func(ctx context.Context, requestInfo map[string]any, request any, handler Handler) (resp any, err error) {
		fmt.Println("limit")
		// 限流逻辑（省略实际实现）
		return handler(ctx, request)
	}
}
