package system_error

type Err struct {
	error
	code    string
	message string
}

// 需要捕获的错误码
var ( // 哨兵错误
	AccessExpire  = NewError("100001", "access token is disabled")  // token已经失效
	RefreshExpire = NewError("100002", "refresh token is disabled") // token已经失效
	SessionExpire = NewError("100003", "session is disabled")       // session已经失效

	TooManyRequests = NewError("100010", "too many requests") // 请求过于频繁
	Unauthorized    = NewError("100011", "unauthorized")      // 没有权限

	CircuitBreaker = NewError("100020", "Circuit Breaker") //熔断器触发
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
