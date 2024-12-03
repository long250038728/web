package server

import "github.com/long250038728/web/tool/register"

const AuthorizationKey = "authorization"
const TraceParentKey = "traceparent"

type Server interface {
	Start() error
	Stop() error
	ServiceInstance() *register.ServiceInstance
}
