package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/server"
	"github.com/long250038728/web/tool/tracing"
	"net/http"
)

var _ server.Server = &Server{}

type HandlerFunc func(*gin.Engine)

type Server struct {
	server      *http.Server
	handler     http.Handler
	address     string
	port        int
	svcInstance *register.ServiceInstance
}

// NewHttp  构造函数
func NewHttp(serverName, address string, port int, handlerFunc HandlerFunc) *Server {
	handler := gin.Default()
	handler.Use(tracing.HandlerFunc())
	svc := &Server{
		server:  &http.Server{Addr: fmt.Sprintf("%s:%d", address, port), Handler: handler},
		handler: handler,
		address: address,
		port:    port,
	}

	svc.svcInstance = &register.ServiceInstance{
		Name:    register.HttpServerName(serverName),
		ID:      register.HttpServerId(serverName),
		Address: svc.address,
		Port:    svc.port,
		Type:    "HTTP",
	}

	handlerFunc(handler)
	return svc
}

// Start 开始启动服务
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Stop 停止服务
func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}

// ServiceInstance 获取服务实例（服务注册与发现）
func (s *Server) ServiceInstance() *register.ServiceInstance {
	return s.svcInstance
}
