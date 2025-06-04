package domain

import (
	"context"
	"github.com/long250038728/web/application/user/internal/repository"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/server/http/gateway"
	"time"
)

type UserDomain struct {
	repository *repository.UserRepository
}

func NewUserDomain(repository *repository.UserRepository) *UserDomain {
	return &UserDomain{repository: repository}
}

func (s *UserDomain) SayHello(ctx context.Context, request *user.RequestHello) (*user.ResponseHello, error) {
	str, err := s.repository.GetName(ctx, request)
	return &user.ResponseHello{Str: str}, err
}
func (s *UserDomain) File(ctx context.Context, request *user.RequestHello) (gateway.FileInterface, error) {
	return &fileDemo{}, nil
}

func (s *UserDomain) SendSSE(ctx context.Context, request *user.RequestHello) (<-chan string, error) {
	c := make(chan string)
	go func() {
		defer func() {
			close(c)
		}()
		for i := 0; i < 100; i++ {
			c <- "hello   "
			time.Sleep(time.Second)
		}
	}()

	return c, nil
}

//=================================================================================================

type fileDemo struct {
	gateway.FileInterface
}

func (f *fileDemo) FileName() string {
	return "file.txt"
}
func (f *fileDemo) FileData() []byte {
	return []byte("this is file data")
}
