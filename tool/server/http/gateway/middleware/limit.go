package middleware

import (
	"context"
	"fmt"
	"github.com/long250038728/web/tool/server/http/gateway"
)

// Limit 示例中间件：限流拦截器
func Limit() gateway.ServerInterceptor {
	return func(ctx context.Context, requestInfo map[string]any, request any, handler gateway.Handler) (resp any, err error) {
		fmt.Println("limit")
		// 限流逻辑（省略实际实现）
		return handler(ctx, request)
	}
}
