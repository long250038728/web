package service

import (
	"context"
	"github.com/long250038728/web/application/agent/internal/domain"
	"github.com/long250038728/web/protoc/agent"
	"github.com/long250038728/web/tool/server/rpc/tool"
)

type Service struct {
	agent.UnimplementedAuthServer
	tool.GrpcHealth
	domain *domain.Domain
}

type AgentServerOpt func(s *Service)

func SetDomain(domain *domain.Domain) AgentServerOpt {
	return func(s *Service) {
		s.domain = domain
	}
}

func NewService(opts ...AgentServerOpt) *Service {
	s := &Service{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Service) Logs(ctx context.Context, req *agent.LogsRequest) (*agent.LogsResponse, error) {
	return s.domain.Logs(ctx, req)
}

func (s *Service) Events(ctx context.Context, req *agent.EventsRequest) (*agent.EventsResponse, error) {
	return s.domain.Events(ctx, req)
}

func (s *Service) Resources(ctx context.Context, req *agent.ResourcesRequest) (*agent.ResourcesResponse, error) {
	return s.domain.Resources(ctx, req)
}
