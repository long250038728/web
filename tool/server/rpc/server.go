package rpc

import (
	"fmt"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/server"
	"github.com/long250038728/web/tool/server/rpc/tool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
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
		grpc.ChainUnaryInterceptor(tool.ServerAuthInterceptor(), tool.ServerTelemetryInterceptor()),
		grpc.KeepaliveEnforcementPolicy(enforcementPolicy), //保持Keepalive
		grpc.KeepaliveParams(serverParameters),             //保持Keepalive
	}
	svc := &Server{
		server:      grpc.NewServer(opts...),
		address:     address,
		port:        port,
		svcInstance: register.NewServiceInstance(serverName, address, port, register.InstanceTypeGRPC),
	}
	fmt.Printf("%s : %s:%d\n", svc.svcInstance.Name, svc.svcInstance.Address, svc.svcInstance.Port)
	handlerFunc(svc.server)

	// grpcurl -plaintext 192.168.0.26:9011 list //查看服务
	// grpcurl -plaintext 192.168.0.26:9011 describe agent.Agent //查看服务中接口列表
	// grpcurl -plaintext -d '{}' 192.168.0.26:9011 agent.Agent/Events //访问接口
	reflection.Register(svc.server) //注册服务反射功能，使客户端可以动态查询服务和方法信息。它在开发和调试过程中极其有用
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
	s.server.GracefulStop()
	return nil
}

// ServiceInstance 获取服务实例（服务注册与发现）
func (s *Server) ServiceInstance() *register.ServiceInstance {
	return s.svcInstance
}
