package consul

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/server/http"
)

type Register struct {
	client  *api.Client
	address string
}

type Config struct {
	Address string `json:"address" yaml:"address"`
}

// NewConsulRegister   创建consul服务注册类
func NewConsulRegister(conf *Config) (*Register, error) {
	//创建consul客户端
	config := api.DefaultConfig()
	config.Address = conf.Address
	config.HttpClient = http.NewCustomHttpClient()
	client, err := api.NewClient(config)

	status := client.Status()
	fmt.Println(status)

	if err != nil {
		return nil, err
	}
	return &Register{client: client, address: conf.Address}, nil
}

// Register 服务注册
func (r *Register) Register(ctx context.Context, serviceInstance *register.ServiceInstance) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	//需要注册的内容
	registration := api.AgentServiceRegistration{
		ID:      serviceInstance.ID,
		Name:    serviceInstance.Name,
		Address: serviceInstance.Address,
		Port:    serviceInstance.Port,
	}

	if serviceInstance.Type == register.InstanceTypeHttp {
		check := api.AgentServiceCheck{}
		check.Timeout = "30s"
		check.Interval = "30s"
		check.DeregisterCriticalServiceAfter = "30s"
		check.HTTP = fmt.Sprintf("http://%s:%d", serviceInstance.Address, serviceInstance.Port)
		registration.Check = &check
	}

	if serviceInstance.Type == register.InstanceTypeGRPC {
		check := api.AgentServiceCheck{}
		check.Timeout = "30s"
		check.Interval = "30s"
		check.DeregisterCriticalServiceAfter = "30s"
		check.GRPC = fmt.Sprintf("%s:%d", serviceInstance.Address, serviceInstance.Port)
		registration.Check = &check
	}

	return r.client.Agent().ServiceRegister(&registration)
}

// DeRegister 服务注销
func (r *Register) DeRegister(ctx context.Context, serviceInstance *register.ServiceInstance) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return r.client.Agent().ServiceDeregister(serviceInstance.ID)
}

// List 获取服务列表
func (r *Register) List(ctx context.Context, serviceName string) ([]*register.ServiceInstance, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

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
	wp, err := watch.Parse(map[string]interface{}{
		"type":    "service",
		"service": serviceName,
	})
	if err != nil {
		return nil, err
	}

	// 定义service变化后所执行的程序(函数)handler
	wp.Handler = func(idx uint64, data interface{}) {
		switch d := data.(type) {
		case []*api.ServiceEntry:
			for _, i := range d {
				// 这里是单个service变化时需要做的逻辑，可以自己添加，或在外部写一个类似handler的函数传进来
				fmt.Printf("service %s 已变化", i.Service.Service)
				// 打印service的状态
				fmt.Println("service status: ", i.Checks.AggregatedStatus())
			}
		}
	}
	// 启动监控
	go wp.Run(r.address)

	return nil, nil
}
