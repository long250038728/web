package server

import "github.com/long250038728/web/tool/register"

type Server interface {
	Start() error
	Stop() error
	ServiceInstance() *register.ServiceInstance
}
