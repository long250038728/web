package http

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"io"
	"net/http"
)

type Middleware struct {
	ginContext *gin.Context
	errorData  map[error]Err

	requestData string
	err         error
}

func NewMiddleware(opts ...MiddlewareOpt) *Middleware {
	errData := make(map[error]Err, 0)
	middle := &Middleware{
		errorData: errData,
	}
	for _, opt := range opts {
		opt(middle)
	}
	return middle
}

func (m *Middleware) GinContext(gin *gin.Context) *Middleware {
	m.ginContext = gin
	return m
}

func (m *Middleware) Bind(data interface{}) *Middleware {
	//读数据
	if m.ginContext.Request.Method == http.MethodGet {
		err := m.ginContext.ShouldBindQuery(data)
		if err != nil {
			m.err = err
		}
		return m
	}

	b, err := io.ReadAll(m.ginContext.Request.Body)
	defer m.ginContext.Request.Body.Close()

	if err != nil {
		m.err = err
		return m
	}

	//如果有数据才读
	if len(b) > 0 {
		err = json.Unmarshal(b, data)
		if err != nil {
			m.err = err
			return m
		}
		m.requestData = string(b)
	} else {
		m.requestData = "{\"request\":\"\"}"
	}

	return m
}

func (m *Middleware) Do(function func(ctx context.Context) (interface{}, error)) {
	ctx := opentelemetry.ExtractHttp(m.ginContext.Request.Context(), m.ginContext.Request)
	span := opentelemetry.NewSpan(ctx, "HTTP Request")
	defer span.Close()
	span.AddEvent(m.requestData)

	if m.err != nil {
		response, _ := m.responseJSON(nil, m.err)
		span.AddEvent(string(response))
		return
	}

	//执行
	data, err := function(span.Context())
	if err != nil {
		span.SetAttributes("err", true)
	}
	response, _ := m.responseJSON(data, err)
	span.AddEvent(string(response))
}

func (m *Middleware) responseJSON(data interface{}, err error) ([]byte, error) {
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
	_, err = m.ginContext.Writer.Write(b)
	return b, err
}

func (m *Middleware) Reset() {
	m.ginContext = nil
	m.err = nil
	m.requestData = ""
}
