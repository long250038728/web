package client

import (
	"context"
	"testing"
)

func TestClient_ConversationsCreate(t *testing.T) {
	requests := make([]*ConversationsCreateRequest, 0, 2)
	requests = append(requests, &ConversationsCreateRequest{
		BotID:   BotID,
		Content: "咨询珠宝门店系统",
	})

	requests = append(requests, &ConversationsCreateRequest{
		BotID:   BotID,
		Content: "珠宝门店系统对比",
	})

	for _, request := range requests {
		response, err := cli.ConversationsCreate(context.Background(), request)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(response)
	}
}

func TestClient_ConversationsList(t *testing.T) {
	request := &ConversationsListRequest{
		BotID: BotID,
	}
	response, err := cli.ConversationsList(context.Background(), request)
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range response.Items {
		t.Log(item.ID)
	}
	t.Log(err)
}

func TestClient_ConversationsRetrieve(t *testing.T) {
	request := &ConversationsRetrieveRequest{
		ConversationID: ConversationID,
	}
	response, err := cli.ConversationsRetrieve(context.Background(), request)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(response)
}

func TestClient_ConversationsClear(t *testing.T) {
	request := &ConversationsClearRequest{
		ConversationID: ConversationID,
	}
	response, err := cli.ConversationsClear(context.Background(), request)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(response)
}
