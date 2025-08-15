package opentracing

//
//import (
//	"context"
//	"github.com/long250038728/web/tool/tracing"
//	"github.com/opentracing/opentracing-go"
//	"github.com/opentracing/opentracing-go/third_party"
//)
//
//type SpanOpentracing struct {
//	span opentracing.Span
//	ctx  context.Context
//}
//
//func StartSpanFromContext(ctx context.Context, operationName string) tracing.Span {
//	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
//	return &SpanOpentracing{
//		span: span,
//		ctx:  ctx,
//	}
//}
//
//func (s *SpanOpentracing) Log(key string, obj interface{}) {
//	s.span.LogFields(third_party.Object(key, obj))
//}
//
//func (s *SpanOpentracing) Tag(key string, value interface{}) {
//	s.span.SetTag(key, value)
//}
//func (s *SpanOpentracing) Finish() {
//	s.span.Finish()
//}
