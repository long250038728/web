package tool

//
//type MiddlewareOpt func(middleware *Middleware)
//
//func SetErrData(data map[error]Err) MiddlewareOpt {
//	return func(middleware *Middleware) {
//		middleware.error = data
//	}
//}
//
//func SetAuth(auth auth.Auth) MiddlewareOpt {
//	return func(middleware *Middleware) {
//		middleware.auth = auth
//	}
//}
//
//func SetLimiter(limiter limiter.Limiter) MiddlewareOpt {
//	return func(middleware *Middleware) {
//		middleware.limiter = limiter
//	}
//}
//
//type response struct {
//	Code    string      `json:"code"`
//	Message string      `json:"message"`
//	Data    interface{} `json:"data"`
//}
//
//type Err struct {
//	Code    string `json:"code"`
//	Message string `json:"message"`
//}
