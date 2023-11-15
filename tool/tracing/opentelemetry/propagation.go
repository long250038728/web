package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

//
//	注入：通过SpanContext注入到carrier_map		—— 把当前的链路相关内容注入到data中（把生成的数据传入下一个服务）
//	opentracing.GlobalTracer().Inject(sm SpanContext, format interface{}, carrier interface{}) error
//		sm: 链路中的spanContext对象
//		format: 类型只有三种   0:Binary   1:TextMap   2:HTTPHeaders
//		carrier: 1跟2 必须实现 TextMapWriter  TextMapReader  这个两个interface{}的实现的对象
//
//	提取：从carrier_map数据中提取生成SpanContext  	  —— 把data中的数据提取出之前的链路中（接收上个服务的数据，转换为本服务的链路）
//	opentracing.GlobalTracer().Extract(format interface{}, carrier interface{}) (SpanContext, error)
//		format:类型只有三种   0:Binary   1:TextMap   2:HTTPHeaders
//		carrier：1跟2 必须实现 TextMapWriter  TextMapReader  这个两个interface{}的实现的对象
//

// Inject 注入 把链路信息注入到map中
func Inject(ctx context.Context, req *http.Request) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
}

// Extract 提取 把map提取信息到context
func Extract(ctx context.Context, req *http.Request) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(req.Header))
}
