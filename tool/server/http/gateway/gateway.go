package gateway

import (
	"context"
	"github.com/gin-gonic/gin"
)

func Json(ginContext *gin.Context, request any) *Gateway {
	return NewGateway(NewJson(ginContext, request))
}

func File(ginContext *gin.Context, request any) *Gateway {
	return NewGateway(NewFile(ginContext, request))
}

func SSE(ginContext *gin.Context, request any) *Gateway {
	return NewGateway(NewSSE(ginContext, request))
}

//========================================================================

// Gateway 拦截器管理器
type Gateway struct {
	Writer
	interceptors ServerInterceptor
	err          error
}

func NewGateway(writer Writer, err error) *Gateway {
	tool := &Gateway{}
	tool.Writer = writer
	tool.err = err
	return tool
}

// Use 添加拦截器链
func (tool *Gateway) Use(interceptors ...ServerInterceptor) *Gateway {
	// 构建拦截器链
	chain := func(ctx context.Context, requestInfo map[string]any, request any, handler Handler) (any, error) {
		for i := len(interceptors) - 1; i >= 0; i-- {
			current := interceptors[i]
			next := handler // 保存当前的 handler，避免闭包捕获
			handler = func(ctx context.Context, request any) (any, error) {
				return current(ctx, requestInfo, request, next)
			}
		}
		return handler(ctx, request)
	}
	tool.interceptors = chain
	return tool
}

// Handle 执行拦截器链和业务逻辑
func (tool *Gateway) Handle(handler Handler) {
	if tool.err != nil {
		tool.WriteErr(tool.err)
		return
	}
	tool.Writer.Run(tool.handle(handler))
}

// Handle 执行拦截器链和业务逻辑
func (tool *Gateway) handle(handler Handler) (any, error) {
	request, requestInfo := tool.Writer.Request()
	if tool.interceptors != nil {
		return tool.interceptors(tool.Writer.Context(), requestInfo, request, handler)
	}
	return handler(tool.Writer.Context(), request)
}
