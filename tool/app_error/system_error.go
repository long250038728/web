package app_error

type Err struct {
	error
	code    string
	message string
}

// 需要捕获的错误码
var ( // 哨兵错误
	ClaimsNull  = NewError("100001", "claims is null")
	SessionNull = NewError("100002", "session is null")

	AccessExpire  = NewError("100003", "access token is disabled")  // token已经失效（重新获取token）
	SessionExpire = NewError("100004", "session is disabled")       // session已经失效 (退出登录)
	RefreshExpire = NewError("100005", "refresh token is disabled") // token已经失效 (退出登录)

	TooManyRequests = NewError("100010", "too many requests") // 请求过于频繁 (http中间件)
	Unauthorized    = NewError("100011", "unauthorized")      // 没有权限
	CircuitBreaker  = NewError("100012", "circuit breaker")   // 熔断器触发
	Validate        = NewError("100013", "validate error")    // 参数校验有误
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
