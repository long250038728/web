package http

import (
	"errors"
	"github.com/long250038728/web/tool/system_error"
)

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(data interface{}, err error) *Response {
	res := &Response{Code: "000000", Message: "success", Data: data}
	if err == nil {
		return res
	}
	var systemErr *system_error.Err
	if ok := errors.As(err, &systemErr); ok == true {
		res.Code = systemErr.Code()
		res.Message = systemErr.Error()
		return res
	}
	res.Code = "999999"
	res.Message = err.Error()
	return res
}
