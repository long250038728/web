package app

import (
	"context"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/server"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Option func(app *App) error

func Servers(servers ...server.Server) Option {
	return func(app *App) error {
		app.servers = servers
		return nil
	}
}

func Register(register register.Register) Option {
	return func(app *App) error {
		app.register = register
		return nil
	}
}

func Tracing(exporter trace.SpanExporter, serviceName string) Option {
	return func(app *App) error {
		t, err := opentelemetry.NewTrace(context.Background(), exporter, serviceName)
		app.trace = t
		return err
	}
}
