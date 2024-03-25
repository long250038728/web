package tool

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type MiddleErr struct {
	Code    string
	Message string
}

func (err *MiddleErr) Error() string {
	return err.Message
}

func NewError(code, message string) *MiddleErr {
	return &MiddleErr{
		Code:    code,
		Message: message,
	}
}
