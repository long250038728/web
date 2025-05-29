package service

import (
	"context"
	"github.com/long250038728/web/application/agent/internal/domain"
	"github.com/long250038728/web/protoc/agent"
	"github.com/long250038728/web/tool/server/rpc/tool"
)

var _ agent.AgentServer = &AgentService{}

type AgentService struct {
	agent.UnimplementedAgentServer
	tool.GrpcHealth
	domain *domain.AgentDomain
}

type AgentServerOpt func(s *AgentService)

func SetDomain(domain *domain.AgentDomain) AgentServerOpt {
	return func(s *AgentService) {
		s.domain = domain
	}
}

func NewService(opts ...AgentServerOpt) *AgentService {
	s := &AgentService{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *AgentService) Logs(ctx context.Context, req *agent.LogsRequest) (*agent.LogsResponse, error) {
	return s.domain.Logs(ctx, req)
}

func (s *AgentService) Events(ctx context.Context, req *agent.EventsRequest) (*agent.EventsResponse, error) {
	return s.domain.Events(ctx, req)
}

func (s *AgentService) Resources(ctx context.Context, req *agent.ResourcesRequest) (*agent.ResourcesResponse, error) {
	return s.domain.Resources(ctx, req)
}
