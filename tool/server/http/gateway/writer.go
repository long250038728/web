package gateway

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

	Write(data interface{})
	WriteErr(err error)
}

// ========================================================================

type writer struct {
	ginContext  *gin.Context
	ctx         context.Context
	request     any
	requestInfo map[string]any
}

func (w *writer) Context() context.Context {
	return w.ctx
}

func (w *writer) Request() (request any, requestInfo map[string]any) {
	return w.request, w.requestInfo
}

func (w *writer) Bind(request any) error {
	if w.ginContext.Request.Method == http.MethodGet {
		return binding.MapFormWithTag(request, w.ginContext.Request.URL.Query(), "json") //get请求需要绑定from标签，查看源码后改为json处理
	} else {
		return w.ginContext.ShouldBind(request)
	}
}

// ========================================================================

func (w *writer) Write(data interface{}) {
	w.WriteJSON(data, nil)
}

func (w *writer) WriteErr(err error) {
	w.WriteJSON(nil, err)
}

func (w *writer) WriteJSON(data interface{}, err error) {
	res := NewResponse(data, err)
	marshalByte, err := json.Marshal(res)
	if err != nil {
		res.Code = "999999"
		res.Message = err.Error()
		res.Details = nil
		marshalByte, _ = json.Marshal(res)
	}
	w.WriteBytes(marshalByte)
}

func (w *writer) WriteBytes(b []byte) {
	//记录请求响应
	w.addLog(string(b))

	w.ginContext.Header("Content-Type", "application/JSON")
	_, _ = w.ginContext.Writer.Write(b)
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

//========================================================================

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

//==========================================================================================================

func newJson(ginContext *gin.Context, request any) (Writer, error) {
	w := &jsonWriter{writer: writer{ginContext: ginContext, request: request}}

	w.ctx = ginContext.Request.Context()
	w.request = request

	// 绑定request参数
	if err := w.Bind(request); err != nil {
		return w, err
	}

	// 把request转为 map[string]any 用于后续处理
	requestInfo, err := w.structToMap(request)
	if err != nil {
		return w, err
	}
	w.requestInfo = requestInfo

	return w, nil
}
func newFile(ginContext *gin.Context, request any) (Writer, error) {
	w := &fileWriter{writer: writer{ginContext: ginContext, request: request}}
	w.ctx = ginContext.Request.Context()
	w.request = request

	// 绑定request参数
	if err := w.Bind(request); err != nil {
		return w, err
	}

	// 把request转为 map[string]any 用于后续处理
	requestInfo, err := w.structToMap(request)
	if err != nil {
		return w, err
	}
	w.requestInfo = requestInfo

	return w, nil
}
func newSse(ginContext *gin.Context, request any) (Writer, error) {
	w := &sseWriter{writer: writer{ginContext: ginContext, request: request}}
	w.ctx = ginContext.Request.Context()
	w.request = request

	// 绑定request参数
	if err := w.Bind(request); err != nil {
		return w, err
	}

	// 把request转为 map[string]any 用于后续处理
	requestInfo, err := w.structToMap(request)
	if err != nil {
		return w, err
	}
	w.requestInfo = requestInfo

	return w, nil
}
