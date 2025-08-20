package tool

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/long250038728/web/tool/app_error"
	"github.com/long250038728/web/tool/authorization"
	"github.com/long250038728/web/tool/server"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ServerTelemetryInterceptor 链路拦截器
func ServerTelemetryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		transparent, ok := md[server.TraceParentKey]
		if !ok || len(transparent) != 1 {
			return handler(ctx, req)
		}

		//写入链路
		ctx = opentelemetry.ExtractMap(ctx, map[string]string{server.TraceParentKey: transparent[0]})
		span := opentelemetry.NewSpan(ctx, "GRPC "+info.FullMethod)
		defer span.Close()

		span.AddEvent(req)
		ctx = span.Context()

		resp, err = handler(ctx, req)

		span.AddEvent(resp)
		if err != nil {
			span.AddEvent(err.Error())
		}
		return resp, err
	}
}

// ServerAuthInterceptor 鉴权拦截器
func ServerAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok && err == nil {
			if claims, ok := md["claims"]; ok && len(claims) == 1 {
				userClaims := &authorization.UserInfo{}
				if err := json.Unmarshal([]byte(claims[0]), &userClaims); err == nil {
					ctx = authorization.SetClaims(ctx, userClaims)
				}
			}
			if session, ok := md["session"]; ok && len(session) == 1 {
				userSession := &authorization.UserSession{}
				if err := json.Unmarshal([]byte(session[0]), &userSession); err == nil {
					ctx = authorization.SetSession(ctx, userSession)
				}
			}
		}
		return handler(ctx, req)
	}
}

// ServerCircuitInterceptor 熔断拦截器（可通过errors.Is(err, app_error.CircuitBreaker)进行判断进行服务降级而不是直接报错）
// 注意: 如果一个函数依赖多个下游，而且他们是有非常密切的关系，中间某个熔断就会导致数据有问题，此时熔断一般是一些不重要的接口，在返回熔断后可用通过缓存或默认值进行服务降级的操作
// 当一个接口依赖了非常重要的下游时，应该控制的是调用接口的入口(服务限流的方式)。以确保整个函数执行过程中不会因为某个下游导致数据不一致问题
func ServerCircuitInterceptor(circuits []string) grpc.UnaryClientInterceptor {
	circuitHash := make(map[string]struct{}, len(circuits))
	for _, circuitPath := range circuits {
		circuitHash[circuitPath] = struct{}{}

		hystrix.ConfigureCommand(circuitPath, hystrix.CommandConfig{
			Timeout:                1000,  // 超时时间（毫秒）
			MaxConcurrentRequests:  10,    // 最大并发请求数（允许最大的请求数，超过熔断）
			RequestVolumeThreshold: 5,     // 触发熔断的最小请求数(当请求多少时开始计算错误率)
			ErrorPercentThreshold:  50,    // 错误百分比阈值 (错误达到多少时触发熔断)
			SleepWindow:            30000, // 熔断器打开后的冷却时间（毫秒）(多久后断路器进入半开状态)
		})
	}
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 如果不在可以熔断降级的接口里面就不错处理
		if _, ok := circuitHash[method]; !ok {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		// 如果是触发熔断的，由于没有真正执行invoker，是没有真正的error的。此时需要返回一个自定义的错误
		err := hystrix.Do(method, func() error { //go get github.com/afex/hystrix-go/hystrix
			return invoker(ctx, method, req, reply, cc, opts...)
		}, nil)
		if err != nil && (errors.Is(err, hystrix.ErrCircuitOpen) || errors.Is(err, hystrix.ErrMaxConcurrency) || errors.Is(err, hystrix.ErrTimeout)) {
			return app_error.CircuitBreaker
		}

		// 真正执行invoker了，是应该返回接口原本的error错误
		return err
	}
}
