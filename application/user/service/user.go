package service

import (
	"context"
	"github.com/long250038728/web/application/user/domain"
	user "github.com/long250038728/web/application/user/protoc"
)

type UserService struct {
	domain *domain.UserDomain
	user.UnimplementedUserServerServer
}

type UserServerOpt func(s *UserService)

func UserDomain(domain *domain.UserDomain) UserServerOpt {
	return func(s *UserService) {
		s.domain = domain
	}
}

func NewUserService(opts ...UserServerOpt) *UserService {
	s := &UserService{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *UserService) SayHello(ctx context.Context, request *user.RequestHello) (*user.ResponseHello, error) {
	return s.domain.SayHello(ctx, request)
}
