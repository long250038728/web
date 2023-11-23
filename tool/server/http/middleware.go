package http

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Middleware struct {
	ginContext *gin.Context
	errorData  map[error]Err
	err        error
}

func NewMiddleware(c *gin.Context, opts ...MiddlewareOpt) *Middleware {
	errData := make(map[error]Err, 0)
	middle := &Middleware{
		ginContext: c,
		errorData:  errData,
	}

	for _, opt := range opts {
		opt(middle)
	}

	return middle
}

func (m *Middleware) Bind(data interface{}) *Middleware {
	//读数据
	b, err := io.ReadAll(m.ginContext.Request.Body)
	defer m.ginContext.Request.Body.Close()

	if err != nil {
		m.err = err
		return m
	}
	//转换数据
	err = json.Unmarshal(b, data)
	if err != nil {
		m.err = err
		return m
	}
	return m
}

func (m *Middleware) response(data interface{}, err error) (int, error) {
	code := "000000"
	message := "success"

	if err != nil {
		code = "999999"
		message = err.Error()
		if errMsg, ok := m.errorData[err]; ok {
			code = errMsg.Code
			message = errMsg.Message
		}
		m.ginContext.Writer.WriteHeader(http.StatusBadRequest)
	}
	res := &response{Code: code, Message: message, Data: data}

	b, _ := json.Marshal(res)
	return m.ginContext.Writer.Write(b)
}

func (m *Middleware) Do(run func(ctx context.Context) (interface{}, error)) {
	if m.err != nil {
		_, _ = m.response(nil, m.err)
		return
	}
	_, _ = m.response(run(context.Background()))
}
