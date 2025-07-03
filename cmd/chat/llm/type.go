package llm

import (
	"errors"
	"github.com/sashabaranov/go-openai"
)

var ModelError = errors.New("model is not empty")

type Config struct {
	BaseURL     string `json:"base_url" yaml:"base_url"`
	AccessToken string `json:"access_token" yaml:"access_token"`
	Model       string `json:"model" yaml:"model"`
}

//====================================================================================

type Conversation interface {
	AddMessage(message string)
	AddResponse(message openai.ChatCompletionMessage)
	List() []openai.ChatCompletionMessage
}
