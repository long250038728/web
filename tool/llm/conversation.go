package llm

import "github.com/sashabaranov/go-openai"

type ConversationMemory struct {
	list []openai.ChatCompletionMessage
}

func NewConversationMemoryLocal(prompt string) Conversation {
	list := make([]openai.ChatCompletionMessage, 0, 100)
	list = append(list, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: prompt,
	})
	return &ConversationMemory{list: list}
}

func (c *ConversationMemory) AddMessage(message string) {
	c.list = append(c.list, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})
}

func (c *ConversationMemory) AddResponse(message openai.ChatCompletionMessage) {
	c.list = append(c.list, message)
}

func (c *ConversationMemory) List() []openai.ChatCompletionMessage {
	return c.list
}
