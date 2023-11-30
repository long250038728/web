package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Span struct {
	spanName string
	context  context.Context
	span     trace.Span
}

func NewSpan(ctx context.Context, spanName string) *Span {
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	return &Span{
		spanName: spanName, context: ctx, span: span,
	}
}

func (s *Span) AddEvent(event string) {
	s.span.AddEvent(event)
}

func (s *Span) SetAttributes(k string, v bool) {
	s.span.SetAttributes(attribute.Bool(k, v))
}

func (s *Span) Context() context.Context {
	return s.context
}

func (s *Span) Close() {
	s.span.End()
}
