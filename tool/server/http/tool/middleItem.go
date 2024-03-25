package tool

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
)

type MiddleItem struct {
	gCtx  *gin.Context
	error map[error]MiddleErr
}

func (e *MiddleItem) Set(gCtx *gin.Context) {
	e.gCtx = gCtx
}

func (e *MiddleItem) SetError(error map[error]MiddleErr) {
	e.error = error
}

func (e *MiddleItem) Ctx() context.Context {
	return opentelemetry.ExtractHttp(e.gCtx.Request.Context(), e.gCtx.Request)
}

func (e *MiddleItem) Reset() {
	e.gCtx = nil
}

func (e *MiddleItem) WriteJSON(data interface{}, err error) {
	ginContext := e.gCtx

	res := e.response(data, err)
	marshalByte, err := json.Marshal(res)

	if err != nil {
		res.Code = "999999"
		res.Message = err.Error()
		res.Data = nil
		marshalByte, _ = json.Marshal(res)
	}

	ginContext.Header("Content-Type", "application/json")
	_, _ = ginContext.Writer.Write(marshalByte)
}

func (e *MiddleItem) response(data interface{}, err error) (res *Response) {
	res = &Response{Code: "000000", Message: "success"}
	if err == nil {
		res.Data = data
		return
	}

	var middleErr *MiddleErr
	if ok := errors.As(err, &middleErr); ok == true {
		res.Code = middleErr.Code
		res.Message = middleErr.Message
		return
	}

	res.Code = "999999"
	res.Message = err.Error()
	return
}
