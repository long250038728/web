package opentelemetry

import "go.opentelemetry.io/otel/exporters/jaeger"

// NewJaegerExporter   http://link.zhubaoe.cn:14268/api/traces
func NewJaegerExporter(endpoint string) (*jaeger.Exporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
}
