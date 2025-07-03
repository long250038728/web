package service

import (
	"github.com/google/wire"
	"github.com/long250038728/web/application/agent/internal/domain"
	"github.com/long250038728/web/protoc/agent"
	"github.com/long250038728/web/tool/server/rpc/tool"
)

var _ agent.AgentServer = &AgentService{}

var ProviderSet = wire.NewSet(NewService)

type AgentService struct {
	agent.UnimplementedAgentServer
	tool.GrpcHealth
	domain *domain.AgentDomain
}

func NewService(domain *domain.AgentDomain) *AgentService {
	s := &AgentService{domain: domain}
	return s
}
