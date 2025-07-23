package opentracing

//
//import (
//	"context"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/long250038728/web/tool/tracing"
//	"github.com/opentracing/opentracing-go/third_party"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/metadata"
//	"strings"
//)
//var Id = "tracing_id"
//
//// HandlerFunc 链路中间件  ———— http
//func HandlerFunc() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		if c.Request.URL.Path == "/favicon.ico" {
//			c.Next()
//			return
//		}
//		span, ctx := extract("HTTP: "+c.Request.URL.Path, c.Request.Context(), c.Request.Header[Id])
//		c.Request = c.Request.WithContext(ctx)
//
//		// 输出响应头, 方便前端调试
//		c.Header(strings.ToUpper(tracing.Id), fmt.Sprintf("%+v", span))
//		c.Next()
//		span.Finish()
//	}
//}
//
//// Interceptor 链路中间件  ———— grpc
//func Interceptor() grpc.UnaryServerInterceptor {
//	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
//		md, _ := metadata.FromIncomingContext(ctx)
//		span, ctx := extract("GRPC: "+info.FullMethod, ctx, md[tracing.Id])
//		span.LogFields(third_party.Object("request", req))
//		resp, err = handler(ctx, req)
//		span.LogFields(third_party.Object("response", resp))
//		span.LogFields(third_party.Object("err", err))
//		span.Finish()
//
//		fmt.Println(resp)
//		return resp, err
//	}
//}
//
//// Context 生成带链路的context  ———— grpc
//func Context(ctx context.Context) context.Context {
//	//处理链路，加到md中
//	carrier, _ := inject(ctx)
//	return metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{tracing.Id: carrier[Id]}))
//}
