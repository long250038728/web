package service

import (
	"context"
	"github.com/long250038728/web/application/auth/internal/domain"
	"github.com/long250038728/web/protoc/auth"
	"github.com/long250038728/web/tool/server/rpc/tool"
)

type Service struct {
	auth.UnimplementedAuthServer
	tool.GrpcHealth
	domain *domain.Domain
}

type UserServerOpt func(s *Service)

func SetDomain(domain *domain.Domain) UserServerOpt {
	return func(s *Service) {
		s.domain = domain
	}
}

func NewService(opts ...UserServerOpt) *Service {
	s := &Service{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Service) Login(ctx context.Context, request *auth.LoginRequest) (*auth.UserResponse, error) {
	return s.domain.Login(ctx, request)
}

func (s *Service) Refresh(ctx context.Context, request *auth.RefreshRequest) (*auth.UserResponse, error) {
	return s.domain.Refresh(ctx, request)
}
