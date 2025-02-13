package llm

import "github.com/sashabaranov/go-openai"

type ConversationMemoryLocal struct {
	list []openai.ChatCompletionMessage
}

func NewConversationMemoryLocal(prompt string) ConversationMemory {
	list := make([]openai.ChatCompletionMessage, 0, 100)
	list = append(list, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: prompt,
	})
	return &ConversationMemoryLocal{list: list}
}

func (c *ConversationMemoryLocal) AddMessage(message string) {
	c.list = append(c.list, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})
}

func (c *ConversationMemoryLocal) AddResponse(message openai.ChatCompletionMessage) {
	c.list = append(c.list, message)
}

func (c *ConversationMemoryLocal) List() []openai.ChatCompletionMessage {
	return c.list
}
