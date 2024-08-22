package tool

import (
	"github.com/long250038728/web/tool/auth/auth"
	"github.com/long250038728/web/tool/limiter"
	"github.com/pkg/errors"
)

type MiddlewareOpt func(middle *Middleware)

// Error 错误列表赋值
func Error(errList []*MiddleErr) MiddlewareOpt {
	return func(middle *Middleware) {
		// O(1)的效率
		var hash = make(map[error]*MiddleErr, len(errList))

		// 通过code快速查找error
		for _, err := range errList {
			hash[errors.New(err.Code)] = err
		}
		//赋值到中间件
		middle.error = hash
	}
}

// Limiter 限速
func Limiter(limiter limiter.Limiter) MiddlewareOpt {
	return func(middle *Middleware) {
		middle.limiter = limiter
	}
}

// Auth 授权
func Auth(auth auth.Auth) MiddlewareOpt {
	return func(middle *Middleware) {
		middle.auth = auth
	}
}
