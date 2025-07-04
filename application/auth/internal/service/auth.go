package service

import (
	"context"
	"github.com/long250038728/web/application/auth/internal/domain"
	"github.com/long250038728/web/protoc/auth"
	"github.com/long250038728/web/tool/server/rpc/tool"
)

type Auth struct {
	auth.UnimplementedAuthServer
	tool.GrpcHealth
	domain *domain.Auth
}

type AuthServerOpt func(s *Auth)

func SetDomain(domain *domain.Auth) AuthServerOpt {
	return func(s *Auth) {
		s.domain = domain
	}
}

func NewService(opts ...AuthServerOpt) *Auth {
	s := &Auth{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Auth) Login(ctx context.Context, request *auth.LoginRequest) (*auth.UserResponse, error) {
	return s.domain.Login(ctx, request)
}

func (s *Auth) Refresh(ctx context.Context, request *auth.RefreshRequest) (*auth.UserResponse, error) {
	return s.domain.Refresh(ctx, request)
}
