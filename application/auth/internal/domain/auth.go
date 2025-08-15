package domain

import (
	"context"
	"github.com/long250038728/web/application/auth/internal/repository"
	"github.com/long250038728/web/protoc/auth"
)

type Auth struct {
	repository *repository.Auth
}

func NewAuthDomain(repository *repository.Auth) *Auth {
	return &Auth{repository: repository}
}

func (s *Auth) Login(ctx context.Context, request *auth.LoginRequest) (*auth.UserResponse, error) {
	return s.repository.Login(ctx, request.Name, request.Password)
}

func (s *Auth) Refresh(ctx context.Context, request *auth.RefreshRequest) (*auth.UserResponse, error) {
	return s.repository.Refresh(ctx, request.RefreshToken)
}
