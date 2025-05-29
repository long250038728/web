package service

import (
	"context"
	"github.com/long250038728/web/application/auth/internal/domain"
	"github.com/long250038728/web/protoc/auth"
	"github.com/long250038728/web/tool/server/rpc/tool"
)

type AuthService struct {
	auth.UnimplementedAuthServer
	tool.GrpcHealth
	domain *domain.AuthDomain
}

type AuthServerOpt func(s *AuthService)

func SetDomain(domain *domain.AuthDomain) AuthServerOpt {
	return func(s *AuthService) {
		s.domain = domain
	}
}

func NewService(opts ...AuthServerOpt) *AuthService {
	s := &AuthService{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *AuthService) Login(ctx context.Context, request *auth.LoginRequest) (*auth.UserResponse, error) {
	return s.domain.Login(ctx, request)
}

func (s *AuthService) Refresh(ctx context.Context, request *auth.RefreshRequest) (*auth.UserResponse, error) {
	return s.domain.Refresh(ctx, request)
}
