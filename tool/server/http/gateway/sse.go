package gateway

import (
	"errors"
	"fmt"
	"net/http"
)

type sseWriter struct {
	writer
}

func (w *sseWriter) Run(response any, err error) {
	// 检查是否支持SSE
	if err := w.checkSupportSEE(); err != nil {
		w.WriteErr(err)
		return
	}
	//处理业务
	if err != nil {
		w.WriteErr(err)
		return
	}

	if ch, ok := response.(<-chan string); !ok {
		w.WriteErr(errors.New("the response is not chan string"))
	} else {
		w.writeSSE(ch)
	}
}

func (w *sseWriter) checkSupportSEE() error {
	if _, ok := w.ginContext.Writer.(http.Flusher); !ok {
		return errors.New("streaming unsupported")
	}
	return nil
}

func (w *sseWriter) writeSSE(ch <-chan string) {
	w.ginContext.Header("Content-Type", "text/event-stream")
	w.ginContext.Header("Cache-Control", "no-store")
	w.ginContext.Header("Connection", "keep-alive")

	ginWriter := w.ginContext.Writer
	f, _ := ginWriter.(http.Flusher)

	for message := range ch {
		_, _ = fmt.Fprintf(ginWriter, "data: %s\n\n", message)
		f.Flush()
	}
	w.addLog("SSE send ok")
}
