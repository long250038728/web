package tool

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type HttpFunc func() (interface{}, error)

type FileInterface interface {
	FileName() string
	FileData() []byte
}

type HttpTools struct {
}

func NewHttpTools() *HttpTools {
	return &HttpTools{}
}

// JSON  json返回
func (m *HttpTools) JSON(gin *gin.Context, request any, function HttpFunc) {
	middleware := NewResponseTools(gin)
	//基础处理 （bind绑定  及链路 处理）
	if err := middleware.Bind(request); err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	//处理业务
	res, err := function()
	middleware.WriteJSON(res, err)
}

// File  File返回
//
//	response 必须实现 FileInterface 接口
func (m *HttpTools) File(gin *gin.Context, request any, function HttpFunc) {
	middleware := NewResponseTools(gin)

	//基础处理 （bind绑定  及链路 处理）
	if err := middleware.Bind(request); err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	//处理业务
	res, err := function()
	if err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	if file, ok := res.(FileInterface); !ok {
		middleware.WriteJSON(nil, errors.New("the file is not interface to FileInterface"))
	} else {
		middleware.WriteFile(file.FileName(), file.FileData())
	}
}

// SSE  SSE返回
//
// response 必须是<-chan string
func (m *HttpTools) SSE(gin *gin.Context, request any, function HttpFunc) {
	middleware := NewResponseTools(gin)
	//基础处理 （bind绑定  及链路 处理）
	if err := middleware.Bind(request); err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	// 检查是否支持SSE
	if err := middleware.CheckSupportSEE(); err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	//处理业务
	res, err := function()
	if err != nil {
		middleware.WriteJSON(nil, err)
		return
	}

	if ch, ok := res.(<-chan string); !ok {
		middleware.WriteJSON(nil, errors.New("the response is not chan string"))
	} else {
		middleware.WriteSSE(ch)
	}
}
