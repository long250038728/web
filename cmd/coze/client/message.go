package client

import (
	"context"
	"github.com/coze-dev/coze-go"
)

// ConversationsMessageCreate 创建消息
func (c *Client) ConversationsMessageCreate(ctx context.Context, request *ConversationsMessageCreateRequest) (*Message, error) {
	Role := coze.MessageRoleUser
	ContentType := coze.MessageContentTypeText
	req := &coze.CreateMessageReq{
		ConversationID: request.ConversationID,
		Content:        request.Content,
		Role:           Role,
		ContentType:    ContentType,
	}
	resp, err := c.getApi().Conversations.Messages.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &Message{Role: string(resp.Role), Type: string(resp.Type), Content: resp.Content, ReasoningContent: resp.ReasoningContent, ContentType: string(resp.ContentType), MetaData: resp.MetaData, ID: resp.ID, ConversationID: resp.ConversationID, SectionID: resp.SectionID, BotID: resp.BotID, ChatID: resp.ChatID, CreatedAt: resp.CreatedAt, UpdatedAt: resp.UpdatedAt}, nil
}

// ConversationsMessageList 查看消息列表
func (c *Client) ConversationsMessageList(ctx context.Context, request *ConversationsMessageListRequest) (*ConversationsMessageListResponse, error) {
	req := &coze.ListConversationsMessagesReq{
		ConversationID: request.ConversationID,
		Limit:          100,
	}
	resp, err := c.getApi().Conversations.Messages.List(ctx, req)
	if err != nil {
		return nil, err
	}

	items := make([]*Message, len(resp.Items()), len(resp.Items()))
	for index, item := range resp.Items() {
		items[len(resp.Items())-index-1] = &Message{
			Role: string(item.Role), Type: string(item.Type), Content: item.Content, ReasoningContent: item.ReasoningContent, ContentType: string(item.ContentType), MetaData: item.MetaData, ID: item.ID, ConversationID: item.ConversationID, SectionID: item.SectionID, BotID: item.BotID, ChatID: item.ChatID, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
		}
	}
	return &ConversationsMessageListResponse{Items: items}, nil
}

// ConversationsMessageRetrieve  查看对话消息详情
func (c *Client) ConversationsMessageRetrieve(ctx context.Context, request *ConversationsMessageRetrieveRequest) (*coze.RetrieveConversationsMessagesResp, error) {
	req := &coze.RetrieveConversationsMessagesReq{
		ConversationID: request.ConversationID,
		MessageID:      request.MessageID,
	}
	resp, err := c.getApi().Conversations.Messages.Retrieve(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
