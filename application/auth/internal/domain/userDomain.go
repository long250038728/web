package domain

import (
	"context"
	"github.com/long250038728/web/application/auth/internal/repository"
	"github.com/long250038728/web/protoc/auth"
)

type AuthDomain struct {
	repository *repository.AuthRepository
}

func NewAuthDomain(repository *repository.AuthRepository) *AuthDomain {
	return &AuthDomain{repository: repository}
}

func (s *AuthDomain) Login(ctx context.Context, request *auth.LoginRequest) (*auth.UserResponse, error) {
	return s.repository.Login(ctx, request.Name, request.Password)
}

func (s *AuthDomain) Refresh(ctx context.Context, request *auth.RefreshRequest) (*auth.UserResponse, error) {
	return s.repository.Refresh(ctx, request.RefreshToken)
}
