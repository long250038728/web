package service

import (
	"context"
	"github.com/long250038728/web/application/user/internal/domain"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/server/http/gateway/encode"
	"github.com/long250038728/web/tool/server/rpc/server"
)

type User struct {
	user.UnimplementedUserServer
	server.GrpcHealth
	domain *domain.User
}

type UserServerOpt func(s *User)

func SetDomain(domain *domain.User) UserServerOpt {
	return func(s *User) {
		s.domain = domain
	}
}

func NewService(opts ...UserServerOpt) *User {
	s := &User{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
func (s *User) SayHello(ctx context.Context, request *user.RequestHello) (*user.ResponseHello, error) {
	return s.domain.SayHello(ctx, request)
}

func (s *User) File(ctx context.Context, request *user.RequestHello) (encode.File, error) {
	return s.domain.File(ctx, request)
}

func (s *User) SendSSE(ctx context.Context, request *user.RequestHello) (<-chan string, error) {
	return s.domain.SendSSE(ctx, request)
}
