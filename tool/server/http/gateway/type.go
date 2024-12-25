package gateway

import (
	"context"
)

// Handler 业务逻辑处理器
type Handler func(ctx context.Context, req any) (any, error)

// ServerInterceptor 中间件拦截器
type ServerInterceptor func(ctx context.Context, requestInfo map[string]any, request any, handler Handler) (resp any, err error)
