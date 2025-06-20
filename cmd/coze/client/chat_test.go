package client

import (
	"context"
	"testing"
)

func TestClient_Chat(t *testing.T) {
	req := &ChatRequest{
		ConversationID: ConversationID,
		BotID:          BotID,
		UserID:         "12345",
		Content:        "铢宝益的特色",
	}

	resp, err := cli.Chat(context.Background(), req)
	if err != nil {
		t.Error(err)
		return
	}

	for _, msg := range resp.Items {
		t.Log(msg.Type)
		t.Log(msg.ContentType)
		t.Log(msg.Content)
		t.Log("================")
	}
}

func TestClient_StreamChat(t *testing.T) {
	req := &ChatRequest{
		ConversationID: ConversationID,
		BotID:          BotID,
		UserID:         "12345",
		Content:        "铢宝益的特色",
	}
	ctx := context.Background()
	ch, err := cli.StreamChat(ctx, req)
	if err != nil {
		t.Error(err)
	}
	for chat := range ch {
		if chat.Err != nil {
			t.Error(chat.Err)
		}
		t.Log(chat.Content)
	}
}

// 当前对话的响应数据
func TestClient_List(t *testing.T) {
	req := &ListRequest{
		ConversationID: ConversationID,
		ChatID:         ChatID,
	}

	resp, err := cli.List(context.Background(), req)
	if err != nil {
		t.Error(err)
		return
	}

	for _, msg := range resp.Items {
		t.Log(msg.Type)
		t.Log(msg.ContentType)
		t.Log(msg.Content)
		t.Log("================")
	}
}

// 当前对话的数据
func TestClient_Retrieve(t *testing.T) {
	req := &RetrieveRequest{
		ConversationID: ConversationID,
		ChatID:         ChatID,
	}

	resp, err := cli.Retrieve(context.Background(), req)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(resp)
}
