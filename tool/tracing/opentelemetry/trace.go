package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type Tracer interface {
	Close(ctx context.Context) error
}

// Trace
// 为什么需要OpenTelemetry：
//
//	1.组件的激增导致可能被拆成十几或更多的微服务，部署在不同的容器中（无法统一管理）
//	2.交互路径的复杂，通过grpc/mq等方式，让一个请求可能跨越多个系统（不在同一个进程内）
//	3.动态性与短暂性，pod的动态创建，动态扩容，动态销毁（导致他们的ip及资源可能会频繁变动）
//	4.把多个日志用链路的方式串联起来（日志连续性）
//
// 能满足的目标：
//
//	1.可以跟踪一个请求经过了多少服务，每个之间的处理耗时，错误出现在哪，每条记录中的层次关系及日志本身
//
// 应用场景：
//
//	1.通过Metrics可以展示有什么请求再某个时间段的突然响应时间变化,点击后进入具体的链路
//	2.通过Tracing中预留tag可快速查询对应的链路数据
//	3.通过Logging，分析整个数据数据
type Trace struct {
	provider *trace.TracerProvider
}

func NewTrace(ctx context.Context, exporter trace.SpanExporter, serviceName string) (*Trace, error) {
	// 链路属性（服务名）
	r, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(serviceName)))
	if err != nil {
		return nil, err
	}

	// 创建provider
	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(r),
	)
	otel.SetTracerProvider(provider)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, // W3C Trace Context（标准跨进程传播）
			propagation.Baggage{}),     // Baggage（上下文附加信息）
	)

	return &Trace{
		provider: provider,
	}, nil
}

func (t *Trace) Close(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}
