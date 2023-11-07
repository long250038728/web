package register

import (
	"context"
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
