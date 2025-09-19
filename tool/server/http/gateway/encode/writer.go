package encode

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"net/http"
)

type Writer interface {
	Bind(request any) error

	Context() context.Context
	Request() (request any, requestInfo map[string]any)

	Run(response any, err error)
	WriteErr(err error)
}

// ========================================================================

type writer struct {
	ginContext  *gin.Context
	ctx         context.Context
	request     any
	requestInfo map[string]any
}

func (w *writer) Init() error {
	w.ctx = w.ginContext.Request.Context()

	// 绑定request参数
	if err := w.Bind(w.request); err != nil {
		return err
	}

	// 把request转为 map[string]any 用于后续处理
	requestInfo, err := w.structToMap(w.request)
	if err != nil {
		return err
	}
	w.requestInfo = requestInfo

	// 日志
	if b, err := json.Marshal(requestInfo); err == nil {
		w.addLog(string(b))
	}
	w.addTags(requestInfo)

	return nil
}

func (w *writer) Bind(request any) error {
	if w.ginContext.Request.Method == http.MethodGet {
		return binding.MapFormWithTag(request, w.ginContext.Request.URL.Query(), "json") //get请求需要绑定from标签，查看源码后改为json处理
	} else {
		return w.ginContext.ShouldBind(request)
	}
}

// ========================================================================

func (w *writer) Context() context.Context {
	return w.ctx
}

func (w *writer) Request() (request any, requestInfo map[string]any) {
	return w.request, w.requestInfo
}

// ========================================================================

func (w *writer) addLog(data string) {
	if span, err := opentelemetry.SpanFromContext(w.ctx); err == nil {
		span.AddEvent(data)
	}
}

func (w *writer) addTags(tags map[string]any) {
	if span, err := opentelemetry.SpanFromContext(w.ctx); err == nil {
		span.SetAttributes(tags)
	}
}

// ========================================================================

func (w *writer) structToMap(request any) (map[string]any, error) {
	var requestInfo map[string]any

	// 将结构体转换为 JSON 字符串
	b, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &requestInfo)
	if err != nil {
		return nil, err
	}
	return requestInfo, err
}

// ========================================================================

func NewJson(ginContext *gin.Context, request any) (Writer, error) {
	w := &jsonWriter{writer: writer{ginContext: ginContext, request: request}}
	if err := w.Init(); err != nil {
		return nil, err
	}
	return w, nil
}

func NewXML(ginContext *gin.Context, request any) (Writer, error) {
	w := &xmlWriter{writer: writer{ginContext: ginContext, request: request}}
	if err := w.Init(); err != nil {
		return nil, err
	}
	return w, nil
}

func NewFile(ginContext *gin.Context, request any) (Writer, error) {
	w := &fileWriter{writer: writer{ginContext: ginContext, request: request}}
	if err := w.Init(); err != nil {
		return nil, err
	}
	return w, nil
}

func NewSSE(ginContext *gin.Context, request any) (Writer, error) {
	w := &sseWriter{writer: writer{ginContext: ginContext, request: request}}
	if err := w.Init(); err != nil {
		return nil, err
	}
	return w, nil
}
