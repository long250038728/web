package domain

import (
	"context"
	"github.com/long250038728/web/application/auth/internal/repository"
	"github.com/long250038728/web/protoc/auth_rpc"
)

type Domain struct {
	repository *repository.Repository
}

func NewDomain(repository *repository.Repository) *Domain {
	return &Domain{
		repository: repository,
	}
}

func (s *Domain) Login(ctx context.Context, request *auth_rpc.LoginRequest) (*auth_rpc.UserResponse, error) {
	return s.repository.Login(ctx, request.Name, request.Password)
}

func (s *Domain) Refresh(ctx context.Context, request *auth_rpc.RefreshRequest) (*auth_rpc.UserResponse, error) {
	return s.repository.Refresh(ctx, request.RefreshToken)
}
