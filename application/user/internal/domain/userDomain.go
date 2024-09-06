package domain

import (
	"context"
	"github.com/long250038728/web/application/user/internal/repository"
	"github.com/long250038728/web/protoc/user"
	"time"
)

type Domain struct {
	repository *repository.Repository
}

func NewDomain(repository *repository.Repository) *Domain {
	return &Domain{
		repository: repository,
	}
}

func (s *Domain) SayHello(ctx context.Context, request *user.RequestHello) (*user.ResponseHello, error) {
	str, err := s.repository.GetName(ctx, request)
	return &user.ResponseHello{Str: str}, err
}

func (s *Domain) SendSSE(ctx context.Context, request *user.RequestHello) (<-chan string, error) {
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
