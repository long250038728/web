package client

import (
	"context"
	"github.com/coze-dev/coze-go"
)

type ConversationsCreateRequest struct {
	BotID    string            `json:"bot_id"`
	Content  string            `json:"content"`
	MetaData map[string]string `json:"meta_data,omitempty"`
}

type ConversationsListRequest struct {
	BotID    string `json:"bot_id"`
	PageNum  int    `json:"page_num,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}

type ConversationsListResponse struct {
	Total int
	Items []*Conversation
}

type ConversationsRetrieveRequest struct {
	ConversationID string `json:"conversation_id"`
}

type ConversationsClearRequest struct {
	ConversationID string `json:"conversation_id"`
}

// ConversationsCreate 创建会话
func (c *Client) ConversationsCreate(ctx context.Context, request *ConversationsCreateRequest) (*Conversation, error) {
	oauth, err := c.GetOAuth()
	if err != nil {
		return nil, err
	}

	req := &coze.CreateConversationsReq{
		BotID: request.BotID,
		Messages: []*coze.Message{
			{Role: coze.MessageRoleUser, Content: request.Content},
		},
	}
	resp, err := coze.NewCozeAPI(coze.NewJWTAuth(oauth, nil), coze.WithBaseURL(coze.CnBaseURL)).Conversations.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &Conversation{ID: resp.ID, CreatedAt: resp.CreatedAt, MetaData: resp.MetaData, LastSectionID: resp.LastSectionID}, nil
}

// ConversationsList 查看会话列表
func (c *Client) ConversationsList(ctx context.Context, request *ConversationsListRequest) (*ConversationsListResponse, error) {
	oauth, err := c.GetOAuth()
	if err != nil {
		return nil, err
	}

	req := &coze.ListConversationsReq{
		BotID:    request.BotID,
		PageNum:  request.PageNum,
		PageSize: request.PageSize,
	}
	resp, err := coze.NewCozeAPI(coze.NewJWTAuth(oauth, nil), coze.WithBaseURL(coze.CnBaseURL)).Conversations.List(ctx, req)
	if err != nil {
		return nil, err
	}

	items := make([]*Conversation, 0, request.PageSize)
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
	oauth, err := c.GetOAuth()
	if err != nil {
		return nil, err
	}

	req := &coze.RetrieveConversationsReq{
		ConversationID: request.ConversationID,
	}
	resp, err := coze.NewCozeAPI(coze.NewJWTAuth(oauth, nil), coze.WithBaseURL(coze.CnBaseURL)).Conversations.Retrieve(ctx, req)
	if err != nil {
		return nil, err
	}
	return &Conversation{ID: resp.ID, CreatedAt: resp.CreatedAt, MetaData: resp.MetaData, LastSectionID: resp.LastSectionID}, nil
}

// ConversationsClear 会话清除
func (c *Client) ConversationsClear(ctx context.Context, request *ConversationsClearRequest) (string, error) {
	oauth, err := c.GetOAuth()
	if err != nil {
		return "", err
	}

	req := &coze.ClearConversationsReq{
		ConversationID: request.ConversationID,
	}
	resp, err := coze.NewCozeAPI(coze.NewJWTAuth(oauth, nil), coze.WithBaseURL(coze.CnBaseURL)).Conversations.Clear(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}
