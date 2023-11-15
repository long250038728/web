package opentracing

import (
	"context"
	"github.com/long250038728/web/tool/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
)

//==================================== 链路初始化  ====================================

// OpentracingGlobalTracer 初始化链路
//
//	address: http://hostname:14268/api/traces
func OpentracingGlobalTracer(address, serviceName string) {
	// 设置上报网关(后续可以改为多模式单选：入文件记录、agent代理记录)
	sender := transport.NewHTTPTransport(address)
	// 初始化 链路系统：1.服务名、2.上报规则、3.记录器
	tracer, _ := jaeger.NewTracer(
		serviceName,
		// 设置采样方案(const表示每个span都记录)
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender, jaeger.ReporterOptions.Logger(jaeger.StdLogger)),

		//自定义key
		jaeger.TracerOptions.CustomHeaderKeys(&jaeger.HeadersConfig{
			TraceContextHeaderName: tracing.Id,
		}),
	)

	// 设置 tracer到全局
	opentracing.SetGlobalTracer(tracer)
}

//==================================== 创建 root span  ====================================
//
//	注入：通过SpanContext把data注入到carrier 		—— 把当前的链路相关内容注入到data中（把生成的数据传入下一个服务）
//	opentracing.GlobalTracer().Inject(sm SpanContext, format interface{}, carrier interface{}) error
//		sm: 链路中的spanContext对象
//		format: 类型只有三种   0:Binary   1:TextMap   2:HTTPHeaders
//		carrier: 1跟2 必须实现 TextMapWriter  TextMapReader  这个两个interface{}的实现的对象
//
//	提取：从carrier数据中提取SpanContext  	 	    —— 把data中的数据提取出之前的链路中（接收上个服务的数据，转换为本服务的链路）
//	opentracing.GlobalTracer().Extract(format interface{}, carrier interface{}) (SpanContext, error)
//		format:类型只有三种   0:Binary   1:TextMap   2:HTTPHeaders
//		carrier：1跟2 必须实现 TextMapWriter  TextMapReader  这个两个interface{}的实现的对象
//

// extract   从carrier数据中提取SpanContext
func extract(operationName string, context context.Context, tracingId []string) (opentracing.Span, context.Context) {
	var opts []opentracing.StartSpanOption
	if parentContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier{tracing.Id: tracingId}); err == nil {
		opts = append(opts, ext.RPCServerOption(parentContext))
	}
	span := opentracing.StartSpan(operationName, opts...)
	return span, opentracing.ContextWithSpan(context, span)
}

// inject  通过SpanContext把data注入到carrier
func inject(context context.Context) (opentracing.TextMapCarrier, error) {
	span := opentracing.SpanFromContext(context)
	if span == nil {
		return opentracing.TextMapCarrier{}, nil
	}
	carrier := opentracing.TextMapCarrier{}
	error := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, carrier)
	return carrier, error
}
