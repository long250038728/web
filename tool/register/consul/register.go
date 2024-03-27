package consul

import (
	"context"
	"github.com/hashicorp/consul/api"
	"github.com/long250038728/web/tool/register"
)

type Register struct {
	client *api.Client
}

type Config struct {
	Address string `json:"address" yaml:"address"`
}

// NewConsulRegister   创建consul服务注册类
func NewConsulRegister(conf *Config) (*Register, error) {
	//创建consul客户端
	config := api.DefaultConfig()
	config.Address = conf.Address
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Register{client: client}, nil
}

// Register 服务注册
func (r *Register) Register(ctx context.Context, serviceInstance *register.ServiceInstance) error {
	//需要注册的内容
	registration := api.AgentServiceRegistration{
		ID:      serviceInstance.ID,
		Name:    serviceInstance.Name,
		Address: serviceInstance.Address,
		Port:    serviceInstance.Port,
	}

	//check := api.AgentServiceCheck{}
	//check.Timeout = "5s"
	//check.Interval = "5s"
	//check.DeregisterCriticalServiceAfter = "30s"
	//check.HTTP = fmt.Sprintf("http://%s:%d", serviceInstance.Address, serviceInstance.Port)
	//registration.Check = &check

	return r.client.Agent().ServiceRegister(&registration)
}

// DeRegister 服务注销
func (r *Register) DeRegister(ctx context.Context, serviceInstance *register.ServiceInstance) error {
	return r.client.Agent().ServiceDeregister(serviceInstance.ID)
}

// List 获取服务列表
func (r *Register) List(ctx context.Context, serviceName string) ([]*register.ServiceInstance, error) {
	// 获取服务列表（可以加缓存）
	svcList, _, err := r.client.Health().Service(serviceName, "", true, nil)

	if err != nil {
		return nil, err
	}

	var list = make([]*register.ServiceInstance, 0, len(svcList))
	for _, svc := range svcList {
		list = append(list, &register.ServiceInstance{
			Name:    serviceName,
			ID:      svc.Service.ID,
			Address: svc.Service.Address,
			Port:    svc.Service.Port,
		})
	}
	return list, nil
}

// Subscribe 监听服务变化
func (r *Register) Subscribe(ctx context.Context, serviceName string) (<-chan *register.ServiceInstance, error) {
	//TODO implement me
	panic("implement me")
}
