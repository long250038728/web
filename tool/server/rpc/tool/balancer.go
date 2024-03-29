package tool

import (
	"github.com/long250038728/web/tool/register"
	"math/rand"
)

type Balancer interface {
	Balancer([]*register.ServiceInstance) *register.ServiceInstance
}

func NewRandBalancer() Balancer {
	return &RandBalancer{}
}

type RandBalancer struct{}

func (b *RandBalancer) Balancer(balancers []*register.ServiceInstance) *register.ServiceInstance {
	return balancers[rand.Int()%len(balancers)]
}
