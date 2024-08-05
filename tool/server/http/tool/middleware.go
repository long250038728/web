package tool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/auth"
	"github.com/long250038728/web/tool/limiter"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"google.golang.org/grpc/metadata"
	"net/http"
	"net/url"
)

type Middleware struct {
	ginContext *gin.Context
	span       *opentelemetry.Span
	ctx        context.Context

	auth    auth.Auth
	limiter limiter.Limiter
	error   map[error]*MiddleErr
}

func NewMiddleware(opts ...MiddlewareOpt) *Middleware {
	middleware := &Middleware{
		error: map[error]*MiddleErr{},
		ctx:   context.Background(),
	}
	for _, opt := range opts {
		opt(middleware)
	}
	return middleware
}

func (m *Middleware) Context(ginContext *gin.Context) (context.Context, error) {
	//从请求里面获取生成新的ctx
	ctx := opentelemetry.ExtractHttp(ginContext.Request.Context(), ginContext.Request)

	//生成一个新的带有span的context
	span := opentelemetry.NewSpan(ctx, ginContext.Request.RequestURI) //记录请求头

	ctx = span.Context()

	//把链路信息生成到metadata中
	mCarrier := map[string]string{}
	mCarrier["authorization"] = ginContext.GetHeader("Authorization")
	opentelemetry.InjectMap(span.Context(), mCarrier) //通过span对象生成map放入metadata中（之后的grpc中获取）

	//把基础信息写到ctx及链路中
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(mCarrier))
	_ = span.Add(mCarrier)
	ginContext.Header("traceparent", mCarrier["traceparent"])

	//限流
	if m.limiter != nil {
		if err := m.limiter.Allow(ctx, ginContext.Request.URL.Path); err != nil {
			return ctx, err
		}
	}

	//授权
	if m.auth != nil {
		parseCtx, err := m.auth.Parse(ctx, ginContext.GetHeader("Authorization"))
		if err != nil {
			return ctx, err
		}
		if err = m.auth.Auth(ctx, ginContext.Request.URL.Path); err != nil {
			return ctx, err
		}
		ctx = parseCtx
	}

	m.ctx = ctx
	m.span = span
	m.ginContext = ginContext

	return m.ctx, nil
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

func (m *Middleware) Reset() {
	if m.ginContext != nil {
		m.ginContext = nil
	}
	if m.span != nil {
		m.span.Close()
		m.span = nil
	}
}

//========================================================================

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

//========================================================================

func (m *Middleware) WriteFile(fileName string, data []byte) {
	fileName = url.QueryEscape(fileName) // 防止中文乱码
	ginContext := m.ginContext
	ginContext.Header("Content-Type", "application/octet-stream")
	ginContext.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	_, _ = ginContext.Writer.Write(data)

	m.span.AddEvent("file is write ok")
}

// ========================================================================

func (m *Middleware) CheckSupportSEE() error {
	if _, ok := m.ginContext.Writer.(http.Flusher); !ok {
		return errors.New("streaming unsupported")
	}
	return nil
}

func (m *Middleware) WriteSSE(ch <-chan string) {
	ginContext := m.ginContext
	ginContext.Header("Content-Type", "text/event-stream")
	ginContext.Header("Cache-Control", "no-cache")
	ginContext.Header("Connection", "keep-alive")

	w := ginContext.Writer
	f, _ := w.(http.Flusher)

	for message := range ch {
		_, _ = fmt.Fprintf(w, message)
		f.Flush()
	}
	m.span.AddEvent("sse send ok")
}

//========================================================================

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
