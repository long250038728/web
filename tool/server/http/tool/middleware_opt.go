package tool

import (
	errors2 "errors"
	"github.com/long250038728/web/tool/auth"
	"github.com/long250038728/web/tool/limiter"
)

type MiddlewareOpt func(middle *Middleware)

func Error(errors []MiddleErr) MiddlewareOpt {
	return func(middle *Middleware) {
		var error = make(map[error]MiddleErr, len(errors))
		for _, err := range errors {
			error[errors2.New(err.Code)] = err
		}
		middle.error = error
	}
}

func Limiter(limiter limiter.Limiter) MiddlewareOpt {
	return func(middle *Middleware) {
		middle.limiter = limiter
	}
}

func Auth(auth auth.Auth) MiddlewareOpt {
	return func(middle *Middleware) {
		middle.auth = auth
	}
}
