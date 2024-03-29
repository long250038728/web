package tool

import (
	errors2 "errors"
	"github.com/long250038728/web/tool/auth"
	"github.com/long250038728/web/tool/limiter"
)

type MiddlewareOpt func(middle *MiddlewarePool)

// Error 错误列表赋值
func Error(errors []MiddleErr) MiddlewareOpt {
	return func(middle *MiddlewarePool) {
		// O(1)的效率
		var hash = make(map[error]MiddleErr, len(errors))

		// 通过code快速查找error
		for _, err := range errors {
			hash[errors2.New(err.Code)] = err
		}

		//赋值到中间件
		middle.error = hash
	}
}

// Limiter 限速
func Limiter(limiter limiter.Limiter) MiddlewareOpt {
	return func(middle *MiddlewarePool) {
		middle.limiter = limiter
	}
}

// Auth 授权
func Auth(auth auth.Auth) MiddlewareOpt {
	return func(middle *MiddlewarePool) {
		middle.auth = auth
	}
}
