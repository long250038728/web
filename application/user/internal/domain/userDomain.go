package domain

import (
	"context"
	"github.com/long250038728/web/application/user/internal/repository"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/server/http/gateway"
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
func (s *Domain) File(ctx context.Context, request *user.RequestHello) (gateway.FileInterface, error) {
	return &fileDemo{}, nil
}

func (s *Domain) SendSSE(ctx context.Context, request *user.RequestHello) (<-chan string, error) {
	//return llm.NewOpenAiClient(llm.SetMessage([]openai.ChatCompletionMessage{
	//	{
	//		Role:    openai.ChatMessageRoleSystem,
	//		Content: `You are a Kubernetes expert. You can write Kubernetes related yaml file.`,
	//	},
	//})).ChatStream(ctx, "i want to deploy a service in kubernetes, i have a docker image is ccr.ccs.tencentyun.com/linl/user:v1 , exposing ports 8001 and 9001")

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
