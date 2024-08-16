package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/server"
	"log"
	"net/http"
	_ "net/http/pprof"
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

	go func() {
		log.Println(http.ListenAndServe("192.168.0.5:6060", nil))
	}()

	handler := gin.Default()
	svc := &Server{
		server:      &http.Server{Addr: fmt.Sprintf("%s:%d", address, port), Handler: handler},
		handler:     handler,
		address:     address,
		port:        port,
		svcInstance: register.NewServiceInstance(serverName, address, port, register.InstanceTypeHttp),
	}

	handler.Use(func(c *gin.Context) {
		// 设置跨域相关的头部信息
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		// 如果是预检请求（OPTIONS 方法），则直接返回成功响应
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatusJSON(http.StatusNoContent, gin.H{}) // 204 No Content
			return
		}

		c.Next()
	})
	fmt.Printf("service %s: %s:%d\n", svc.svcInstance.Type, svc.svcInstance.Address, svc.svcInstance.Port)

	//health check
	handler.GET("/health", func(gin *gin.Context) {
		gin.Writer.WriteHeader(http.StatusOK)
		_, _ = gin.Writer.Write([]byte("health"))
	})

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
