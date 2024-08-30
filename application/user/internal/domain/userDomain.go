package domain

import (
	"context"
	"github.com/long250038728/web/application/user/internal/repository"
	"github.com/long250038728/web/protoc/user"
	"time"
)

type UserDomain struct {
	userRepository *repository.UserRepository
}

func NewDomain(userRepository *repository.UserRepository) *UserDomain {
	return &UserDomain{
		userRepository: userRepository,
	}
}

func (s *UserDomain) SayHello(ctx context.Context, request *user.RequestHello) (*user.ResponseHello, error) {
	str, err := s.userRepository.GetName(ctx, request)
	return &user.ResponseHello{Str: str}, err
}

func (s *UserDomain) SendSSE(ctx context.Context, request *user.RequestHello) (<-chan string, error) {
	ch := make(chan string, 10)
	go func() {
		defer close(ch)
		for i := 0; i < 100; i++ {
			select {
			case <-ctx.Done():
				return
			default:
			}

			ch <- time.Now().Local().Format(time.DateTime) + "\n"
			time.Sleep(time.Second)
		}
	}()
	return ch, nil
}
