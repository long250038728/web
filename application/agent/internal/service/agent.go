package service

import (
	"github.com/google/wire"
	"github.com/long250038728/web/application/agent/internal/domain"
	"github.com/long250038728/web/protoc/agent"
	"github.com/long250038728/web/tool/server/rpc/tool"
)

var _ agent.AgentServer = &Agent{}

var ProviderSet = wire.NewSet(NewService)

type Agent struct {
	agent.UnimplementedAgentServer
	tool.GrpcHealth
	domain *domain.Agent
}

func NewService(domain *domain.Agent) *Agent {
	s := &Agent{domain: domain}
	return s
}
