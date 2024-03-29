package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

// InjectHttp 注入 把链路信息注入到http中
func InjectHttp(ctx context.Context, req *http.Request) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
}

// ExtractHttp 提取 把http提取信息到context
func ExtractHttp(ctx context.Context, req *http.Request) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(req.Header))
}

// InjectMap 注入 把链路信息注入到map中
func InjectMap(ctx context.Context, mCarrier map[string]string) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(mCarrier))
}

// ExtractMap 提取 把map提取信息到context
func ExtractMap(ctx context.Context, mCarrier map[string]string) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(mCarrier))
}
