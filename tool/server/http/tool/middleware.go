package tool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/authorization/session"
	"github.com/long250038728/web/tool/limiter"
	"github.com/long250038728/web/tool/system_error"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"google.golang.org/grpc/metadata"
	"net/http"
	"net/url"
)

type Middleware struct {
	ginContext *gin.Context
	span       *opentelemetry.Span

	auth    session.Auth
	limiter limiter.Limiter
}

func NewMiddleware(opts ...MiddlewareOpt) *Middleware {
	middleware := &Middleware{}
	for _, opt := range opts {
		opt(middleware)
	}
	return middleware
}

func (m *Middleware) Context(ginContext *gin.Context) (context.Context, error) {
	m.ginContext = ginContext
	authorization := ginContext.GetHeader("Authorization")

	//生成一个新的带有span的context
	// 1. 以ginContext.Request.Context() 获取ctx上下文 (也可以通过context.Background()创建)
	// 2. 通过 telemetry 提取 http请求头中的参数生成一个名称为请求URI的 span (如果请求头中有traceparent 则生成一个子span，如果无则生成一个root span)
	// 3. 通过span 获取新的 ctx 以后续使用
	span := opentelemetry.NewSpan(opentelemetry.ExtractHttp(ginContext.Request.Context(), ginContext.Request), ginContext.Request.RequestURI) //记录请求头
	ctx := span.Context()
	m.span = span

	mCarrier := map[string]string{"authorization": authorization} // mCarrier["authorization"] = authorization // 把 http 请求头中的Authorization信息写入mCarrier
	opentelemetry.InjectMap(ctx, mCarrier)                        // 把 telemetry的id等信息写入mCarrier

	_ = span.Add(mCarrier)
	ginContext.Header("traceparent", mCarrier["traceparent"])

	//限流(有authorization信息则以人为单位)
	if m.limiter != nil {
		if err := m.limiter.Allow(ctx, "http:"+authorization); err != nil {
			return ctx, err
		}
	}

	//授权
	if m.auth != nil {
		parseCtx, err := m.auth.Parse(ctx, authorization)
		if err != nil {
			return ctx, err
		}
		if err = m.auth.Auth(ctx, ginContext.Request.URL.Path); err != nil {
			return ctx, err
		}
		ctx = parseCtx
	}

	//把所有信息写入metadata中并生成新的ctx
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(mCarrier))
	return ctx, nil
}

func (m *Middleware) Bind(request any) error {
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

	//这里是Marshal 报错
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
	res = &Response{Code: "000000", Message: "success", Data: data}
	if err == nil {
		return
	}
	var middleErr *system_error.Err
	if ok := errors.As(err, &middleErr); ok == true {
		res.Code = middleErr.Code()
		res.Message = middleErr.Error()
		return
	}

	res.Code = "999999"
	res.Message = err.Error()
	return
}
