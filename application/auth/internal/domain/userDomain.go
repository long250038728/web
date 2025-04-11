package domain

import (
	"context"
	"github.com/long250038728/web/application/auth/internal/repository"
	auth "github.com/long250038728/web/protoc/auth"
)

type Domain struct {
	repository *repository.Repository
}

func NewDomain(repository *repository.Repository) *Domain {
	return &Domain{
		repository: repository,
	}
}

func (s *Domain) Login(ctx context.Context, request *auth.LoginRequest) (*auth.UserResponse, error) {
	return s.repository.Login(ctx, request.Name, request.Password)
}

func (s *Domain) Refresh(ctx context.Context, request *auth.RefreshRequest) (*auth.UserResponse, error) {
	return s.repository.Refresh(ctx, request.RefreshToken)
}
