package llm

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

//go get github.com/sashabaranov/go-openai

type Chat struct {
	conf               *Config
	client             *openai.Client
	conversationMemory Conversation
}

func NewChat(conf *Config, conversationMemory Conversation) (*Chat, error) {
	if conf.Model == "" {
		return nil, ModelError
	}
	config := openai.DefaultConfig(conf.AccessToken)
	config.BaseURL = conf.BaseURL
	return &Chat{
		conf:               conf,
		client:             openai.NewClientWithConfig(config),
		conversationMemory: conversationMemory,
	}, nil
}

func (chat *Chat) Chat(ctx context.Context, message string) (string, error) {
	if chat.conversationMemory == nil {
		chat.conversationMemory = NewConversationMemoryLocal("")
	}

	chat.conversationMemory.AddMessage(message)

	resp, err := chat.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    chat.conf.Model,
			Messages: chat.conversationMemory.List(),
		},
	)

	chat.conversationMemory.AddResponse(resp.Choices[0].Message)

	return resp.Choices[0].Message.Content, err
}
