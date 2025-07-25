package middleware

import (
	"context"
	"github.com/long250038728/web/tool/app_error"
	"github.com/long250038728/web/tool/server/http/gateway"
	"reflect"
)

// Validate 校验请求参数
func Validate(keys []string) gateway.ServerInterceptor {
	return func(ctx context.Context, requestInfo map[string]any, request any, handler gateway.Handler) (resp any, err error) {
		for _, k := range keys {
			d, ok := requestInfo[k]
			if !ok || reflect.DeepEqual(d, reflect.Zero(reflect.TypeOf(d)).Interface()) {
				return nil, app_error.Validate
			}
		}
		return handler(ctx, request)
	}
}
