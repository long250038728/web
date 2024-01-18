go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/sdk
go get go.opentelemetry.io/otel/trace
go.opentelemetry.io/otel/exporters/jaeger

```
    jaeger链路追踪（包含ui，client等）
    oltp类似一个agent代理（jagger自己也有个agent，但是被废弃）
    
    使用opentelemetry有以下方式
    otel库 =>  jaeger                              //otel库可以生成jaeger能接收的格式
    otel库 =>  otel collector  => jaeger/other     //otel库把数据先发到otel collector，代理再把数据转发到jaeger/other中
```



```
docker pull jaegertracing/all-in-one
docker pull otel/opentelemetry-collector

// 16686 Jaeger UI 端口
// 14268 端口是 Jaeger Collector 默认的 HTTP 端口，用于接收追踪数据。（http 端口）
// 6831 端口用于接收来自服务的 Zipkin Thrift 格式的追踪数据。（Compact 端口）
// 6832 端口用于接收来自服务的 Jaeger Thrift 格式的追踪数据。（Compact 端口）
//
docker run -d --name jaeger \
  -p 16686:16686 \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 14268:14268 \
  jaegertracing/all-in-one
  
 
// 4317 端口是 opentelemetry Collector 默认的 HTTP 端口，用于接收追踪数据。 
docker run -d --name otel-collector \
  -v $(pwd)/otel-collector-config.yaml:/etc/otel-collector-config.yaml \
  -p 4317:4317 \
  otel/opentelemetry-collector --config /etc/otel-collector-config.yaml
```



```otel-collector-config.yaml
receivers:
  otlp:  # 定义接收器，使用 OTLP 协议接收数据
    protocols:
      grpc: {}  # 使用 gRPC 协议
      http: {}  # 使用 http 协议

processors:
  batch:  # 定义处理器，将数据进行批处理
    timeout: 5s  # 批处理的超时时间

exporters:
  jaeger:  # 定义导出器，将数据导出到 Jaeger 后端
    endpoint: "http://jaeger:14268/api/traces"  # Jaeger 后端的地址

service:
  pipelines:
    traces:
      receivers: [otlp]  # 使用 otlp 接收器
      processors: [batch]  # 使用 batch 处理器
      exporters: [jaeger]  # 使用 jaeger 导出器

```




```
exporter, err := jaeger.New(jaeger.WithCollectorEndpoint("http://jaeger:14268/api/traces"))    //jaeger collector 
exporter, err := otlp.NewExporter(otlp.WithInsecure(),otlp.WithAddress("otel-collector:4317")) //otel collector

resources, err := resource.Merge(
    resource.Default(),
    resource.NewWithAttributes(semconv.ServiceNameKey.String("go-example-service"))
)

//创建链路提供者   
//      exporter  导出器（导出到哪里）
//      resources 资源（设置服务名等信息）
//      sampler   采样 (设置怎样的采样率)
tracerProvider := sdktrace.NewTracerProvider(
    sdktrace.WithBatcher(exporter),
    sdktrace.WithResource(resources),
    sdktrace.WithSampler(sdktrace.AlwaysSample()),
)

//设置链路提供者为全局（单例）
otel.SetTracerProvider(tracerProvider)

//设置传播类 （用于 Extract && Inject）
otel.SetTextMapPropagator(
	propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{}
	),
)

//链路提供者不用时需要关闭
tracerProvider.Shutdown(context.Background())	
```



```
//从rquest请求头中获取,并生成到ctx中
ctx := otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(req.Header))


//通过ctx新增span,ctx  ctx代表父级中的span
ctx, span := otel.Tracer("").Start(ctx, spanName, opts...)


//把链路信息保存到map中用于传递到其他系统中
mCarrier := make(map[string]string)
otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(mCarrier))
```