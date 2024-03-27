package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"testing"
)

func TestTracing(t *testing.T) {
	ctx := context.Background()

	exporter, err := NewJaegerExporterAddress("http://link.zhubaoe.cn:14268/api/traces")
	if err != nil {
		t.Error(err)
	}

	trace, err := NewTrace(ctx, exporter, "ServiceA")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = trace.Close(ctx)
		t.Error(err)
	}()

	ctx = context.Background()
	span := NewSpan(ctx, "this is first")
	defer span.Close()

	span2 := NewSpan(span.Context(), "this is two")
	defer span2.Close()

	span.AddEvent("hello world")
	span.SetAttributes("hello1", true)

	span2.AddEvent("hello world2")
	span2.SetAttributes("hello2", false)

	//链路注入到map
	mCarrier := make(map[string]string)
	otel.GetTextMapPropagator().Inject(span.Context(), propagation.MapCarrier(mCarrier))

	//map提取到链路
	newCtx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(mCarrier))

	span3 := NewSpan(newCtx, "this is three")
	defer span3.Close()

	span3.AddEvent("hello world3")
	span3.SetAttributes("hello3", false)
}
