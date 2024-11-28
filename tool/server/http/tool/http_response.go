package tool

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"net/http"
	"net/url"
)

type ResponseTools struct {
	g    *gin.Context
	span *opentelemetry.Span
}

func NewResponseTools(ginContext *gin.Context) *ResponseTools {
	middleware := &ResponseTools{g: ginContext}
	if span, err := opentelemetry.SpanFromContext(ginContext.Request.Context()); err == nil {
		middleware.span = span
	}
	return middleware
}

func (m *ResponseTools) Bind(request any) error {
	var err error
	if m.g.Request.Method == http.MethodGet {
		err = m.g.BindQuery(request)
	} else {
		err = m.g.ShouldBind(request)
	}

	//绑定成功记录链路
	if err == nil {
		requestByte, requestErr := json.Marshal(request)
		if requestErr != nil {
			err = requestErr
		}
		//记录请求参数
		m.writeLog(string(requestByte))

		//写入tag
		var requestData map[string]any
		if requestErr = json.Unmarshal(requestByte, &requestData); requestErr == nil {
			for k, v := range requestData {
				m.writeTag(k, v)
			}
		}
	}

	//记录错误信息
	if err != nil {
		m.writeLog(err.Error())
	}
	return err
}

//========================================================================

func (m *ResponseTools) WriteJSON(data interface{}, err error) {
	res := m.response(data, err)
	marshalByte, err := json.Marshal(res)

	//记录请求响应
	m.writeLog(string(marshalByte))

	//这里是Marshal 报错
	if err != nil {
		res.Code = "999999"
		res.Message = err.Error()
		res.Data = nil
		marshalByte, _ = json.Marshal(res)
	}

	m.g.Header("Content-Type", "application/json")
	_, _ = m.g.Writer.Write(marshalByte)
}

//========================================================================

func (m *ResponseTools) WriteFile(fileName string, data []byte) {
	fileName = url.QueryEscape(fileName) // 防止中文乱码
	m.g.Header("Content-Type", "application/octet-stream")
	m.g.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")

	_, _ = m.g.Writer.Write(data)
	m.writeLog("file is write ok")
}

// ========================================================================

func (m *ResponseTools) CheckSupportSEE() error {
	if _, ok := m.g.Writer.(http.Flusher); !ok {
		return errors.New("streaming unsupported")
	}
	return nil
}

func (m *ResponseTools) WriteSSE(ch <-chan string) {
	m.g.Header("Content-Type", "text/event-stream")
	m.g.Header("Cache-Control", "no-cache")
	m.g.Header("Connection", "keep-alive")

	w := m.g.Writer
	f, _ := w.(http.Flusher)

	for message := range ch {
		_, _ = fmt.Fprintf(w, message)
		f.Flush()
	}
	m.writeLog("sse send ok")
}

//========================================================================

func (m *ResponseTools) response(data interface{}, err error) *Response {
	return NewResponse(data, err)
}

func (m *ResponseTools) writeLog(data string) {
	if m.span == nil {
		return
	}
	m.span.AddEvent(data)
}

func (m *ResponseTools) writeTag(key string, value any) {
	if m.span == nil {
		return
	}
	// 只记录部分类型数据
	switch value.(type) {
	case string, byte, []byte,
		float64, float32,
		int, int64, int32, int16, int8,
		uint, uint64, uint32, uint16,
		bool:
		m.span.SetAttributes(key, value)
	default:

	}
}
