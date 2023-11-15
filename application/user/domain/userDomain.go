package domain

import (
	"context"
	user "github.com/long250038728/web/application/user/protoc"
	"github.com/long250038728/web/application/user/repository"
	"github.com/long250038728/web/tool/tracing/opentracing"
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
	span := opentracing.StartSpanFromContext(ctx, "s_domain")
	defer span.Finish()
	span.Log("hello", "name")
	return &user.ResponseHello{Str: "HELLO :" + request.Name}, nil
}
