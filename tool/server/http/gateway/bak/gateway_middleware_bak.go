package bak

//
//import (
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/gin-gonic/gin/binding"
//	"github.com/long250038728/web/tool/tracing/opentelemetry"
//	"net/http"
//	"net/url"
//)
//
//type GatewayMiddleware struct {
//	g    *gin.Context
//	span *opentelemetry.Span
//}
//
//func NewGatewayMiddleware() *GatewayMiddleware {
//	return &GatewayMiddleware{}
//}
//
//func (m *GatewayMiddleware) SetGin(ginContext *gin.Context) *GatewayMiddleware {
//	m.g = ginContext
//	if span, err := opentelemetry.SpanFromContext(ginContext.Request.Context()); err == nil {
//		m.span = span
//	}
//	return m
//}
//
//func (m *GatewayMiddleware) Bind(request any) error {
//	var err error
//	if m.g.Request.Method == http.MethodGet {
//		err = binding.MapFormWithTag(request, m.g.Request.URL.Query(), "json") //get请求需要绑定from标签，查看源码后改为json处理
//	} else {
//		err = m.g.ShouldBind(request)
//	}
//
//	//绑定成功记录链路
//	if err == nil {
//		requestByte, requestErr := json.Marshal(request)
//		if requestErr != nil {
//			err = requestErr
//		}
//		//记录请求参数
//		m.writeLog(string(requestByte))
//
//		//写入tag
//		var requestData map[string]any
//		if requestErr = json.Unmarshal(requestByte, &requestData); requestErr == nil {
//			m.writeTags(requestData)
//		}
//	}
//
//	//记录错误信息
//	if err != nil {
//		m.writeLog(err.Error())
//	}
//	return err
//}
//
//func (m *GatewayMiddleware) Reset() {
//	m.g = nil
//	m.span = nil
//}
//
//// ========================================================================
//
//func (m *GatewayMiddleware) WriteData(data interface{}) {
//	m.WriteJSON(data, nil)
//}
//
//func (m *GatewayMiddleware) WriteErr(err error) {
//	m.WriteJSON(nil, err)
//}
//
//func (m *GatewayMiddleware) WriteJSON(data interface{}, err error) {
//	res := NewResponse(data, nil)
//	marshalByte, err := json.Marshal(res)
//	if err != nil {
//		res.Code = "999999"
//		res.Message = err.Error()
//		res.Data = nil
//		marshalByte, _ = json.Marshal(res)
//	}
//	m.WriteBytes(marshalByte)
//}
//
//func (m *GatewayMiddleware) WriteBytes(b []byte) {
//	//记录请求响应
//	m.writeLog(string(b))
//
//	m.g.Header("Content-Type", "application/JSON")
//	_, _ = m.g.Writer.Write(b)
//}
//
////========================================================================
//
//func (m *GatewayMiddleware) WriteFile(fileName string, data []byte) {
//	fileName = url.QueryEscape(fileName) // 防止中文乱码
//	m.g.Header("Content-Type", "application/octet-stream")
//	m.g.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
//
//	_, _ = m.g.Writer.Write(data)
//	m.writeLog("File is write ok")
//}
//
//// ========================================================================
//
//func (m *GatewayMiddleware) CheckSupportSEE() error {
//	if _, ok := m.g.Writer.(http.Flusher); !ok {
//		return errors.New("streaming unsupported")
//	}
//	return nil
//}
//
//func (m *GatewayMiddleware) WriteSSE(ch <-chan string) {
//	m.g.Header("Content-Type", "text/event-stream")
//	m.g.Header("Cache-Control", "no-store")
//	m.g.Header("Connection", "keep-alive")
//
//	w := m.g.Writer
//	f, _ := w.(http.Flusher)
//
//	for message := range ch {
//		_, _ = fmt.Fprintf(w, message)
//		f.Flush()
//	}
//	m.writeLog("SSE send ok")
//}
//
////========================================================================
//
//func (m *GatewayMiddleware) writeLog(data string) {
//	if m.span == nil {
//		return
//	}
//	m.span.AddEvent(data)
//}
//
//func (m *GatewayMiddleware) writeTags(tags map[string]any) {
//	if m.span == nil {
//		return
//	}
//	m.span.SetAttributes(tags)
//}
