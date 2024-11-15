package tool

import (
	"github.com/long250038728/web/tool/authorization/session"
	"github.com/long250038728/web/tool/limiter"
)

type MiddlewareOpt func(middle *Middleware)

// Limiter 限速
func Limiter(limiter limiter.Limiter) MiddlewareOpt {
	return func(middle *Middleware) {
		middle.limiter = limiter
	}
}

// Auth 授权
func Auth(auth session.Auth) MiddlewareOpt {
	return func(middle *Middleware) {
		middle.auth = auth
	}
}
