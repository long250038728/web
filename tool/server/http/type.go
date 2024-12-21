package http

import (
	"errors"
	"github.com/long250038728/web/tool/app_error"
)

type Response struct {
	Code    string      `JSON:"code"`
	Message string      `JSON:"message"`
	Data    interface{} `JSON:"data"`
}

func NewResponse(data interface{}, err error) *Response {
	res := &Response{Code: "000000", Message: "success", Data: data}
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
