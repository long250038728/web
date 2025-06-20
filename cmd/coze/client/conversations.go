package client

import (
	"context"
	"github.com/coze-dev/coze-go"
)

// ConversationsCreate 创建会话
func (c *Client) ConversationsCreate(ctx context.Context, request *ConversationsCreateRequest) (*Conversation, error) {
	req := &coze.CreateConversationsReq{
		BotID: request.BotID,
		Messages: []*coze.Message{
			{Role: coze.MessageRoleUser, Content: request.Content},
		},
	}
	resp, err := c.getApi().Conversations.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &Conversation{ID: resp.ID, CreatedAt: resp.CreatedAt, MetaData: resp.MetaData, LastSectionID: resp.LastSectionID}, nil
}

// ConversationsList 查看会话列表
func (c *Client) ConversationsList(ctx context.Context, request *ConversationsListRequest) (*ConversationsListResponse, error) {
	req := &coze.ListConversationsReq{
		BotID:    request.BotID,
		PageNum:  request.PageNum,
		PageSize: request.PageSize,
	}
	resp, err := c.getApi().Conversations.List(ctx, req)
	if err != nil {
		return nil, err
	}

	items := make([]*Conversation, 0, len(resp.Items()))
	for _, item := range resp.Items() {
		items = append(items, &Conversation{
			ID:            item.ID,
			CreatedAt:     item.CreatedAt,
			MetaData:      item.MetaData,
			LastSectionID: item.LastSectionID,
		})
	}
	return &ConversationsListResponse{Total: resp.Total(), Items: items}, nil
}

// ConversationsRetrieve 查看会话
func (c *Client) ConversationsRetrieve(ctx context.Context, request *ConversationsRetrieveRequest) (*Conversation, error) {
	req := &coze.RetrieveConversationsReq{
		ConversationID: request.ConversationID,
	}
	resp, err := c.getApi().Conversations.Retrieve(ctx, req)
	if err != nil {
		return nil, err
	}
	return &Conversation{ID: resp.ID, CreatedAt: resp.CreatedAt, MetaData: resp.MetaData, LastSectionID: resp.LastSectionID}, nil
}

// ConversationsClear 会话清除
func (c *Client) ConversationsClear(ctx context.Context, request *ConversationsClearRequest) (string, error) {
	req := &coze.ClearConversationsReq{
		ConversationID: request.ConversationID,
	}
	resp, err := c.getApi().Conversations.Clear(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}
