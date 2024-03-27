package opentelemetry

import "go.opentelemetry.io/otel/exporters/jaeger"

type Config struct {
	Address string `json:"address" yaml:"address"`
}

// NewJaegerExporterAddress   http://link.zhubaoe.cn:14268/api/traces

func NewJaegerExporter(conf *Config) (*jaeger.Exporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(conf.Address)))
}

func NewJaegerExporterAddress(endpoint string) (*jaeger.Exporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
}
