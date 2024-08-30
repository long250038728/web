package domain

import (
	"context"
	"github.com/long250038728/web/application/auth/internal/repository"
	"github.com/long250038728/web/protoc/auth_rpc"
)

type Domain struct {
	userRepository *repository.Repository
}

func NewDomain(userRepository *repository.Repository) *Domain {
	return &Domain{
		userRepository: userRepository,
	}
}

func (s *Domain) Login(ctx context.Context, request *auth_rpc.LoginRequest) (*auth_rpc.UserResponse, error) {
	return s.userRepository.Login(ctx, request.Name, request.Password)
}

func (s *Domain) Refresh(ctx context.Context, request *auth_rpc.RefreshRequest) (*auth_rpc.UserResponse, error) {
	return s.userRepository.Refresh(ctx, request.RefreshToken)
}
