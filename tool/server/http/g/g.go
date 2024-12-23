package g

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UnaryHandler 业务逻辑处理器
type UnaryHandler func(ctx context.Context, req map[string]any) (any, error)

// UnaryServerInterceptor 中间件拦截器
type UnaryServerInterceptor func(ctx context.Context, req map[string]any, handler UnaryHandler) (resp any, err error)

// AAA 拦截器管理器
type AAA struct {
	interceptor UnaryServerInterceptor
	ctx         context.Context
	request     map[string]any
}

func NewAAA(g *gin.Context, request any) *AAA {
	var err error
	if g.Request.Method == http.MethodGet {
		err = g.BindJSON(request)
	} else {
		err = g.ShouldBind(request)
	}
	fmt.Println(err, request)

	return &AAA{ctx: g.Request.Context(), request: map[string]any{}}
}

// Use 添加拦截器链
func (a *AAA) Use(interceptors ...UnaryServerInterceptor) *AAA {
	// 构建拦截器链
	chain := func(ctx context.Context, req map[string]any, handler UnaryHandler) (any, error) {
		for i := len(interceptors) - 1; i >= 0; i-- {
			current := interceptors[i]
			next := handler // 保存当前的 handler，避免闭包捕获
			handler = func(ctx context.Context, req map[string]any) (any, error) {
				return current(ctx, req, next)
			}
		}
		return handler(ctx, req)
	}
	a.interceptor = chain
	return a
}

// Handle 执行拦截器链和业务逻辑
func (a *AAA) Handle(handler UnaryHandler) (any, error) {
	if a.interceptor != nil {
		return a.interceptor(a.ctx, a.request, handler)
	}
	return handler(a.ctx, a.request)
}

// 示例中间件：缓存拦截器
func cache() UnaryServerInterceptor {
	return func(ctx context.Context, req map[string]any, handler UnaryHandler) (any, error) {
		fmt.Println("cache")
		// 检查缓存逻辑（省略实际实现）
		return handler(ctx, req)
	}
}

// 示例中间件：限流拦截器
func limit() UnaryServerInterceptor {
	return func(ctx context.Context, req map[string]any, handler UnaryHandler) (any, error) {
		fmt.Println("limit")
		// 限流逻辑（省略实际实现）
		return handler(ctx, req)
	}
}
