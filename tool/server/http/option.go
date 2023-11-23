package http

type MiddlewareOpt func(middleware *Middleware)

func SetErrData(data map[error]Err) MiddlewareOpt {
	return func(middleware *Middleware) {
		middleware.errorData = data
	}
}
