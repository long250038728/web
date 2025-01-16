package llm

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"testing"
)

func TestChat(t *testing.T) {
	chat := NewOpenAiClient()
	t.Log(chat.Chat(context.Background(), "1+100=? Just give me a number result"))
}

func TestChatStream(t *testing.T) {
	chat := NewOpenAiClient(SetMessage([]openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: `You are a Kubernetes expert. You can write Kubernetes related yaml file.`,
		},
	}))
	ch, err := chat.ChatStream(context.Background(), "i want to deploy a service in kubernetes, i have a docker image is ccr.ccs.tencentyun.com/linl/user:v1 , exposing ports 8001 and 9001,把所以的yaml文件整理成一个，请用中文回复输出详细讲解")
	if err != nil {
		t.Error(err)
		return
	}

	bytes := make([]byte, 0, 0)
	for str := range ch {
		bytes = append(bytes, []byte(str)...)
	}
	t.Log(string(bytes))
}
