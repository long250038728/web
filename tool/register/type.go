package register

import (
	"context"
	"fmt"
	"math/rand"
)

type Register interface {
	Register(ctx context.Context, serviceInstance *ServiceInstance) error
	DeRegister(ctx context.Context, serviceInstance *ServiceInstance) error
	List(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	Subscribe(ctx context.Context, serviceName string) (<-chan *ServiceInstance, error)
}

type ServiceInstance struct {
	ID      string
	Name    string
	Address string
	Port    int
	Type    string
}

func HttpServerName(serverName string) string {
	return fmt.Sprintf("%v-%v", serverName, "HTTP")
}

func HttpServerId(serverName string) string {
	return fmt.Sprintf("%v-%v-%d", serverName, "HTTP", rand.Uint64()%10000)
}

func GrpcServerName(serverName string) string {
	return fmt.Sprintf("%v-%v", serverName, "GRPC")
}

func GrpcServerId(serverName string) string {
	return fmt.Sprintf("%v-%v-%d", serverName, "GRPC", rand.Uint64()%10000)
}
