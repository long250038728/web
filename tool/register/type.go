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

const InstanceTypeHttp = "HTTP"
const InstanceTypeGRPC = "GRPC"

func NewServiceInstance(serverName, address string, port int, instanceType string) *ServiceInstance {
	instance := &ServiceInstance{
		Address: address,
		Port:    port,
		Type:    instanceType,
	}
	instance.Name = GetServerName(serverName, instanceType)
	instance.ID = GetServerId(serverName, instanceType)
	return instance
}

func GetServerName(serverName string, instanceType string) string {
	return fmt.Sprintf("%v-%v", serverName, instanceType)
}

func GetServerId(serverName string, instanceType string) string {
	return fmt.Sprintf("%v-%v-%d", serverName, instanceType, rand.Uint64()%10000)
}
