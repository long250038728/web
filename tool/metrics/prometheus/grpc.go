package prometheus

import (
	"context"
	"google.golang.org/grpc"
)

// Interceptor  prometheus统计相关
func Interceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		return resp, err
	}
}
