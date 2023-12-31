package domain

import (
	"context"
	user "github.com/long250038728/web/application/user/protoc"
	"github.com/long250038728/web/application/user/repository"
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
