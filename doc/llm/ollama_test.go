package llm

import (
	"context"
	"testing"
)

func TestQwChat_Chat(t *testing.T) {
	chat := NewOpenAiClient()
	t.Log(chat.Chat(context.Background(), "1+100=? Just give me a number result"))
}
