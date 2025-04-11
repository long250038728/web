package domain

import (
	"context"
	"github.com/long250038728/web/application/user/internal/repository"
	"github.com/long250038728/web/protoc/user"
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
	//return llm.NewOpenAiClient(llm.SetMessage([]openai.ChatCompletionMessage{
	//	{
	//		Role:    openai.ChatMessageRoleSystem,
	//		Content: `You are a Kubernetes expert. You can write Kubernetes related yaml file.`,
	//	},
	//})).ChatStream(ctx, "i want to deploy a service in kubernetes, i have a docker image is ccr.ccs.tencentyun.com/linl/user:v1 , exposing ports 8001 and 9001")

	return nil, nil
}
