package app

import (
	"context"
	"github.com/long250038728/web/tool/register/consul"
	"github.com/long250038728/web/tool/server"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"go.opentelemetry.io/otel/exporters/jaeger"
)

type Option func(app *App) error

func Servers(servers ...server.Server) Option {
	return func(app *App) error {
		app.servers = servers
		return nil
	}
}

func Register(register *consul.Register) Option {
	return func(app *App) error {
		if register != nil {
			app.register = register
		}
		return nil
	}
}

func Tracing(exporter *jaeger.Exporter, serviceName string) Option {
	return func(app *App) error {
		if exporter != nil {
			t, err := opentelemetry.NewTrace(context.Background(), exporter, serviceName)
			app.trace = t
			return err
		}
		return nil
	}
}
