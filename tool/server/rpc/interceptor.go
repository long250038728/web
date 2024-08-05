package rpc

import (
	"context"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/auth"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func serverInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		var span *opentelemetry.Span

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			//写入用户信息
			if authorization, ok := md["authorization"]; ok && len(authorization) == 1 {
				ctx, _ = auth.NewAuth(app.NewUtil().Cache()).Parse(ctx, authorization[0])
			}

			//写入链路
			if transparent, ok := md["traceparent"]; ok && len(transparent) == 1 {
				ctx = opentelemetry.ExtractMap(ctx, map[string]string{"traceparent": transparent[0]})
				span = opentelemetry.NewSpan(ctx, "GRPC: "+info.FullMethod)
				defer span.Close()

				_ = span.Add(req)
				ctx = span.Context()

				defer func() {
					span.Add(resp)
					if err != nil {
						span.Add(err.Error())
					}
				}()

			}
		}

		resp, err = handler(ctx, req)
		return resp, err
	}
}
