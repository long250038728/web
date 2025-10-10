package opentelemetry

import (
	"context"
	"encoding/json"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Span struct {
	spanName string
	context  context.Context
	span     trace.Span
}

// SpanFromContext 在ctx中获取span对象
func SpanFromContext(ctx context.Context) (*Span, error) {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return nil, errors.New("span is not valid")
	}
	return &Span{spanName: "", context: ctx, span: span}, nil
}

// NewSpan 创建span对象
func NewSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) *Span {
	// 创建span 当未设置otel.SetTracerProvider(provider)时，返回的span是一个nonRecordingSpan，该对象所有的方法都是空实现，以便于后续无需判断span是否存在的场景
	//   在otel.Tracer("")方法中判断设置TracerProvider用delegate字段表示
	ctx, span := otel.Tracer("").Start(ctx, spanName, opts...)
	return &Span{
		spanName: spanName, context: ctx, span: span,
	}
}

func (s *Span) TraceID() string {
	sContext := s.span.SpanContext()
	return sContext.TraceID().String()
}

func (s *Span) AddEvent(event any) {
	deepNum := 1
	paths := make([]attribute.KeyValue, 0, deepNum)
	//for i := 1; i <= deepNum; i++ {
	//	pc, _, line, ok := runtime.Caller(i)
	//	if !ok {
	//		break
	//	}
	//	//通过 pc 获取函数名
	//	funcName := runtime.FuncForPC(pc).Name()
	//	paths = append(paths, attribute.String("path", fmt.Sprintf("%s:%d", funcName, line)))
	//}

	switch event.(type) {
	case string:
		s.span.AddEvent(event.(string), trace.WithAttributes(paths...))
	default:
		b, _ := json.Marshal(event)
		s.span.AddEvent(string(b), trace.WithAttributes(paths...))
	}
}

func (s *Span) SetAttribute(k string, v any) {
	s.SetAttributes(map[string]any{k: v})
}

func (s *Span) SetAttributes(attributes map[string]any) {
	var kvs = make([]attribute.KeyValue, 0, len(attributes))
	for k, v := range attributes {
		switch v.(type) {
		case string:
			kvs = append(kvs, attribute.String(k, v.(string)))
		case int:
			kvs = append(kvs, attribute.Int(k, v.(int)))
		case int64:
			kvs = append(kvs, attribute.Int64(k, v.(int64)))
		case float64:
			kvs = append(kvs, attribute.Float64(k, v.(float64)))
		case bool:
			kvs = append(kvs, attribute.Bool(k, v.(bool)))
		default:
		}
	}
	if len(kvs) == 0 {
		return
	}
	s.span.SetAttributes(kvs...)
}

func (s *Span) Context() context.Context {
	return s.context
}

func (s *Span) Close() {
	s.span.End()
}
