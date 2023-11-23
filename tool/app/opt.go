package app

import (
	"github.com/long250038728/web/tool/auth"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/server"
)

type Option func(app *App)

func Servers(servers ...server.Server) Option {
	return func(app *App) {
		app.servers = servers
	}
}

func Auth(auth auth.Auth) Option {
	return func(app *App) {
		app.auth = auth
	}
}

func Register(register register.Register) Option {
	return func(app *App) {
		app.register = register
	}
}

func Tracing(serverName, address string) Option {
	return func(app *App) {
		//opentracing.OpentracingGlobalTracer(address, serverName)
	}
}
