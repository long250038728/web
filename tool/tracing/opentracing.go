package tracing

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type SpanOpentracing struct {
	span opentracing.Span
	ctx  context.Context
}

func StartSpanFromContext(ctx context.Context, operationName string) Span {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	return &SpanOpentracing{
		span: span,
		ctx:  ctx,
	}
}

func (s *SpanOpentracing) Log(key string, obj interface{}) {
	s.span.LogFields(log.Object(key, obj))
}

func (s *SpanOpentracing) Tag(key string, value interface{}) {
	s.span.SetTag(key, value)
}
func (s *SpanOpentracing) Finish() {
	s.span.Finish()
}
