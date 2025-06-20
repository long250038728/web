package client

import (
	"context"
	"testing"
)

func TestClient_ConversationsMessageCreate(t *testing.T) {
	request := &ConversationsMessageCreateRequest{
		ConversationID: ConversationID,
		Content:        "今天天气怎么样",
	}
	response, err := cli.ConversationsMessageCreate(context.Background(), request)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(response)
}

func TestClient_ConversationsMessageList(t *testing.T) {
	request := &ConversationsMessageListRequest{
		ConversationID: ConversationID,
	}
	response, err := cli.ConversationsMessageList(context.Background(), request)
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range response.Items {
		t.Log(item)
	}
}

func TestClient_ConversationsMessageRetrieve(t *testing.T) {
	request := &ConversationsMessageRetrieveRequest{
		ConversationID: ConversationID,
		MessageID:      MessageID,
	}
	response, err := cli.ConversationsMessageRetrieve(context.Background(), request)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(response)
}
