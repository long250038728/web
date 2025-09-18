package domain

import (
	"context"
	"github.com/long250038728/web/application/user/internal/repository"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/server/http/gateway/encode"
	"time"
)

type User struct {
	repository *repository.User
}

func NewUserDomain(repository *repository.User) *User {
	return &User{repository: repository}
}

func (s *User) SayHello(ctx context.Context, request *user.RequestHello) (*user.ResponseHello, error) {
	str, err := s.repository.GetName(ctx, request)
	return &user.ResponseHello{Str: str}, err
}
func (s *User) File(ctx context.Context, request *user.RequestHello) (encode.File, error) {
	return &fileDemo{}, nil
}

func (s *User) SendSSE(ctx context.Context, request *user.RequestHello) (<-chan string, error) {
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
	encode.File
}

func (f *fileDemo) FileName() string {
	return "file.txt"
}
func (f *fileDemo) FileData() []byte {
	return []byte("this is file data")
}
