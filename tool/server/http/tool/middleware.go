package tool

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"google.golang.org/grpc/metadata"
	"net/http"
)

type Middleware struct {
	ginContext *gin.Context
	error      map[error]MiddleErr
	span       *opentelemetry.Span
	ctx        context.Context
}

func NewMiddleware(error map[error]MiddleErr) *Middleware {
	return &Middleware{
		error: error,
	}
}

func (m *Middleware) Set(ginContext *gin.Context) {
	m.ginContext = ginContext

	m.ctx = opentelemetry.ExtractHttp(m.ginContext.Request.Context(), ginContext.Request) //从请求里面获取生成新的ctx
	m.span = opentelemetry.NewSpan(m.ctx, "HTTP Request")                                 //生成一个新的带有span的context

	//把链路信息生成到metadata中
	mCarrier := map[string]string{}
	opentelemetry.InjectMap(m.span.Context(), mCarrier) //通过span对象生成map放入metadata中（之后的grpc中获取）
	m.ctx = metadata.NewOutgoingContext(m.span.Context(), metadata.New(mCarrier))

	//记录请求头
	m.span.AddEvent(m.ginContext.Request.RequestURI)
}

func (m *Middleware) bind(request any) error {
	var err error
	if m.ginContext.Request.Method == http.MethodGet {
		err = m.ginContext.BindQuery(request)
	} else {
		err = m.ginContext.ShouldBind(request)
	}

	//绑定成功记录链路
	if err == nil {
		requestByte, requestErr := json.Marshal(request)
		if requestErr != nil {
			err = requestErr
		}
		//记录请求参数
		m.span.AddEvent(string(requestByte))
	}

	//记录错误信息
	if err != nil {
		m.span.AddEvent(err.Error())
	}
	return err
}

func (m *Middleware) Context() context.Context {
	return m.ctx
}

func (m *Middleware) Reset() {
	if m.ginContext != nil {
		m.ginContext = nil
	}
	if m.span != nil {
		m.span.Close()
	}
}

func (m *Middleware) WriteJSON(data interface{}, err error) {
	ginContext := m.ginContext
	res := m.response(data, err)
	marshalByte, err := json.Marshal(res)

	//记录请求响应
	m.span.AddEvent(string(marshalByte))

	if err != nil {
		res.Code = "999999"
		res.Message = err.Error()
		res.Data = nil
		marshalByte, _ = json.Marshal(res)
	}

	ginContext.Header("Content-Type", "application/json")
	_, _ = ginContext.Writer.Write(marshalByte)
}

func (m *Middleware) response(data interface{}, err error) (res *Response) {
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
