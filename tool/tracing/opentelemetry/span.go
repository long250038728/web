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

func SpanFromContext(ctx context.Context) (*Span, error) {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return nil, errors.New("span is not valid")
	}
	return &Span{spanName: "", context: ctx, span: span}, nil
}

func NewSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) *Span {
	ctx, span := otel.Tracer("").Start(ctx, spanName, opts...)
	return &Span{
		spanName: spanName, context: ctx, span: span,
	}
}

func (s *Span) TraceID() string {
	sContext := s.span.SpanContext()
	return sContext.TraceID().String()
}

func (s *Span) Add(event any) error {
	switch event.(type) {
	case string:
		s.AddEvent(event.(string))
	default:
		b, err := json.Marshal(event)
		if err != nil {
			return err
		}
		s.AddEvent(string(b))
	}

	return nil
}

func (s *Span) AddEvent(event string) {
	s.span.AddEvent(event)
}

func (s *Span) SetAttributes(k string, v any) {
	var kvs = make([]attribute.KeyValue, 0, 1)

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
