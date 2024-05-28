package domain

import (
	"context"
	"github.com/long250038728/web/application/user/internal/repository"
	"github.com/long250038728/web/protoc/user"
)

type UserDomain struct {
	userRepository *repository.UserRepository
}

func NewUserDomain(userRepository *repository.UserRepository) *UserDomain {
	return &UserDomain{
		userRepository: userRepository,
	}
}

func (s *UserDomain) SayHello(ctx context.Context, request *user.RequestHello) (*user.ResponseHello, error) {
	str, err := s.userRepository.GetName(ctx, request)
	return &user.ResponseHello{Str: str}, err
}
