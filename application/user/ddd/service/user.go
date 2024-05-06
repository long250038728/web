package service

import (
	"context"
	"github.com/long250038728/web/application/user/ddd/domain"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/server/rpc"
)

type UserService struct {
	user.UnimplementedUserServer
	rpc.GrpcHealth
	domain *domain.UserDomain
}

type UserServerOpt func(s *UserService)

func SetUserDomain(domain *domain.UserDomain) UserServerOpt {
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
