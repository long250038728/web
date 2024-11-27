package rpc

import (
	"context"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/authorization/session"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// serverTelemetryInterceptor 链路拦截器
func serverTelemetryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		transparent, ok := md["traceparent"]
		if !ok || len(transparent) != 1 {
			return handler(ctx, req)
		}

		//写入链路
		ctx = opentelemetry.ExtractMap(ctx, map[string]string{"traceparent": transparent[0]})
		span := opentelemetry.NewSpan(ctx, "GRPC: "+info.FullMethod)
		defer span.Close()

		_ = span.Add(req)
		ctx = span.Context()

		resp, err = handler(ctx, req)

		_ = span.Add(resp)
		if err != nil {
			span.Add(err.Error())
		}
		return resp, err
	}
}

// serverAuthInterceptor 鉴权拦截器
func serverAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok && err == nil {
			cache, err := app.NewUtil().Cache()
			if err != nil {
				return handler(ctx, req)
			}

			if authorization, ok := md["authorization"]; ok && len(authorization) == 1 {
				ctx, _ = session.NewAuth(cache).Parse(ctx, authorization[0])
			}
		}
		return handler(ctx, req)
	}
}

//func serverInterceptor() grpc.UnaryServerInterceptor {
//	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
//		var span *opentelemetry.Span
//		cache, err := app.NewUtil().Cache()
//
//		//接收grpc的md数据
//		if md, ok := metadata.FromIncomingContext(ctx); ok && err == nil {
//			//写入用户信息
//			if authorization, ok := md["authorization"]; ok && len(authorization) == 1 {
//				ctx, _ = session.NewAuth(cache).Parse(ctx, authorization[0])
//			}
//
//			//写入链路
//			if transparent, ok := md["traceparent"]; ok && len(transparent) == 1 {
//				ctx = opentelemetry.ExtractMap(ctx, map[string]string{"traceparent": transparent[0]})
//				span = opentelemetry.NewSpan(ctx, "GRPC: "+info.FullMethod)
//				defer span.Close()
//
//				_ = span.Add(req)
//				ctx = span.Context()
//
//				defer func() {
//					span.Add(resp)
//					if err != nil {
//						span.Add(err.Error())
//					}
//				}()
//
//			}
//		}
//
//		resp, err = handler(ctx, req)
//		return resp, err
//	}
//}
