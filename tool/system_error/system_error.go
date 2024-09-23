package system_error

type Err struct {
	error
	code    string
	message string
}

var ( // 哨兵错误
	AccessExpire  = NewError("100001", "access token is disabled")  // token已经失效
	RefreshExpire = NewError("100002", "refresh token is disabled") // token已经失效
	SessionExpire = NewError("100001", "session is disabled")       // token已经失效

	LimiterTime = NewError("110001", "The limiter time setting is incorrect") // limiter时间设置有误
	Limiter     = NewError("110002", "the limiter has been triggered")        // limiter已经触发
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
