package app_error

type Err struct {
	error
	code    string
	message string
}

// 需要捕获的错误码
var ( // 哨兵错误
	AccessExpire  = NewError("100001", "access token is disabled")  // token已经失效（重新获取token）
	SessionExpire = NewError("100002", "session is disabled")       // session已经失效 (退出登录)
	RefreshExpire = NewError("100003", "refresh token is disabled") // token已经失效 (退出登录)

	TooManyRequests = NewError("100010", "too many requests") // 请求过于频繁 (http中间件)
	Unauthorized    = NewError("100011", "unauthorized")      // 没有权限

	CircuitBreaker = NewError("100020", "circuit breaker") //熔断器触发

	ApiTooManyRequests = NewError("100012", "api too many requests") // 请求过于频繁(接口)
)

func NewError(code, message string) error {
	return &Err{code: code, message: message}
}

func (err *Err) Code() string {
	return err.code
}

func (err *Err) Error() string {
	return err.message
}
