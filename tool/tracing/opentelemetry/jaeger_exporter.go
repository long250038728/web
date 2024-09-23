package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	Address string `json:"address" yaml:"address"`
}

type SpanExporter interface {
	ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error
	Shutdown(ctx context.Context) error
}

// NewJaegerExporter 设置
// Address   http://xxxx.xxxx.com/api/traces
func NewJaegerExporter(conf *Config) (SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(conf.Address)))
}
