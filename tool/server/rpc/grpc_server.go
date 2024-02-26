package rpc

import (
	"fmt"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/server"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

var _ server.Server = &Server{}

type HandlerFunc func(engine *grpc.Server)

// Server 服务端
type Server struct {
	server      *grpc.Server
	handler     http.Handler
	address     string
	port        int
	svcInstance *register.ServiceInstance
}

// NewGrpc  构造函数
func NewGrpc(serverName, address string, port int, handlerFunc HandlerFunc) *Server {
	opts := []grpc.ServerOption{
		//grpc.ChainUnaryInterceptor(opentracing.Interceptor(), prometheus.Interceptor()),
	}

	svc := &Server{
		server:  grpc.NewServer(opts...),
		address: address,
		port:    port,
	}

	svc.svcInstance = &register.ServiceInstance{
		Name:    register.GrpcServerName(serverName),
		ID:      register.GrpcServerId(serverName),
		Address: svc.address,
		Port:    svc.port,
		Type:    "GRPC",
	}
	fmt.Printf("service %s: %s:%d\n", svc.svcInstance.Type, svc.svcInstance.Address, svc.svcInstance.Port)
	handlerFunc(svc.server)
	return svc
}

// Start 开始启动服务
func (s *Server) Start() error {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.address, s.port))
	if err != nil {
		return err
	}
	return s.server.Serve(listen)
}

// Stop 停止服务
func (s *Server) Stop() error {
	s.server.Stop()
	return nil
}

// ServiceInstance 获取服务实例（服务注册与发现）
func (s *Server) ServiceInstance() *register.ServiceInstance {
	return s.svcInstance
}
