package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type Trace struct {
	provider *trace.TracerProvider
}

func NewTrace(context context.Context, exporter trace.SpanExporter, serviceName string) (*Trace, error) {
	// 链路属性（服务名）
	r, err := resource.New(context, resource.WithAttributes(semconv.ServiceName(serviceName)))
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
			propagation.TraceContext{},
			propagation.Baggage{}),
	)

	return &Trace{
		provider: provider,
	}, nil
}

func (t *Trace) Close(context context.Context) error {
	return t.provider.Shutdown(context)
}
