package rpc

import (
	"fmt"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"net/http"
	"time"
)

var _ server.Server = &Server{}
var enforcementPolicy = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

var serverParameters = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
}

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
		//grpc.ChainUnaryInterceptor(opentracing.Interceptor(), prometheus.Interceptor()), //中间件
		grpc.KeepaliveEnforcementPolicy(enforcementPolicy), //保持Keepalive
		grpc.KeepaliveParams(serverParameters),             //保持Keepalive
	}

	svc := &Server{
		server:      grpc.NewServer(opts...),
		address:     address,
		port:        port,
		svcInstance: register.NewServiceInstance(serverName, address, port, "GRPC"),
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
