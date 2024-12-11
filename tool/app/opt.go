package app

import (
	"context"
	"fmt"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/server"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"reflect"
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
		if !IsNil(register) {
			app.register = register
		}
		return nil
	}
}

func Tracing(exporter opentelemetry.SpanExporter, serviceName string) Option {
	return func(app *App) error {
		if !IsNil(exporter) {
			t, err := opentelemetry.NewTrace(context.Background(), exporter, fmt.Sprintf("server-%s", serviceName))
			app.trace = t
			return err
		}
		return nil
	}
}

// IsNil 判断interface是否为空
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}

	vi := reflect.ValueOf(i)

	//判断是否是指针，可通过指针判断指向的内存是不是为空
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	//如果不是指针就代表一定有一个具体的值
	return false
}
