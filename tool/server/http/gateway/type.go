package gateway

import (
	"context"
	"errors"
	"github.com/long250038728/web/tool/app_error"
)

// Handler 业务逻辑处理器
type Handler func(ctx context.Context, req any) (any, error)

// ServerInterceptor 中间件拦截器
type ServerInterceptor func(ctx context.Context, requestInfo map[string]any, request any, handler Handler) (resp any, err error)

// Response 符合swagger返回格式
type Response struct {
	Code    string      `JSON:"code"`
	Message string      `JSON:"message"`
	Details interface{} `JSON:"details"`
}

func NewResponse(data interface{}, err error) *Response {
	res := &Response{Code: "000000", Message: "success", Details: data}
	if err == nil {
		return res
	}
	var systemErr *app_error.Err
	if ok := errors.As(err, &systemErr); ok == true {
		res.Code = systemErr.Code()
		res.Message = systemErr.Error()
		return res
	}
	res.Code = "999999"
	res.Message = err.Error()
	return res
}
