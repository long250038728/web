package encode

import (
	"errors"
	"github.com/long250038728/web/tool/app_error"
)

type File interface {
	FileName() string
	FileData() []byte
}

//====================================================================

// Response 符合swagger返回格式
type Response struct {
	Code    string      `JSON:"code"`
	Message string      `JSON:"message"`
	Details interface{} `JSON:"details"`
}

func NewResponse(data interface{}, err error) *Response {
	res := &Response{Code: "000000", Message: "success", Details: data}
	if err == nil {
		return res
	}
	var systemErr *app_error.Err
	if ok := errors.As(err, &systemErr); ok == true {
		res.Code = systemErr.Code()
		res.Message = systemErr.Error()
		return res
	}
	res.Code = "999999"
	res.Message = err.Error()
	return res
}
