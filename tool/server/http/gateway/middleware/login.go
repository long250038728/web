package middleware

import (
	"context"
	"github.com/long250038728/web/tool/authorization"
	"github.com/long250038728/web/tool/server/http/gateway"
)

// Login 校验登录
func Login() gateway.ServerInterceptor {
	return func(ctx context.Context, requestInfo map[string]any, request any, handler gateway.Handler) (resp any, err error) {
		if _, err = authorization.GetClaims(ctx); err != nil {
			return nil, err
		}
		return handler(ctx, request)
	}
}
