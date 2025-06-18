package client

import (
	"context"
	"errors"
	"github.com/coze-dev/coze-go"
	"io"
	"net/http"
	"time"
)

type ChatRequest struct {
	ConversationID string `json:"conversation_id"`
	BotID          string `json:"bot_id"`
	UserID         string `json:"user_id"`

	Content  string            `json:"content"`
	MetaData map[string]string `json:"meta_data"`
}
type ChatResponse struct {
	Items []*Message
}

type RetrieveRequest struct {
	ConversationId string `json:"conversation_id"`
	ChatID         string `json:"chat_id"`
}

type ListRequest struct {
	ConversationId string `json:"conversation_id"`
	ChatID         string `json:"chat_id"`
}
type ListResponse struct {
	Items []*Message
}

// Chat 创建会话
func (c *Client) Chat(ctx context.Context, request *ChatRequest) (*ListResponse, error) {
	oauth, err := c.GetOAuth()
	if err != nil {
		return nil, err
	}

	req := &coze.CreateChatsReq{
		ConversationID: request.ConversationID,
		BotID:          request.BotID,
		UserID:         request.UserID,
		Messages: []*coze.Message{
			coze.BuildUserQuestionText(request.Content, nil),
		},
	}

	timeout := int(time.Second) * 20
	resp, err := coze.NewCozeAPI(coze.NewJWTAuth(oauth, nil), coze.WithBaseURL(coze.CnBaseURL)).Chat.CreateAndPoll(ctx, req, &timeout)
	if err != nil {
		return nil, err
	}

	items := make([]*Message, 0, len(resp.Messages))
	for _, item := range resp.Messages {
		items = append(items, &Message{
			Role: string(item.Role), Type: string(item.Type), Content: item.Content, ReasoningContent: item.ReasoningContent, ContentType: string(item.ContentType), MetaData: item.MetaData, ID: item.ID, ConversationID: item.ConversationID, SectionID: item.SectionID, BotID: item.BotID, ChatID: item.ChatID, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
		})
	}
	return &ListResponse{items}, nil
}

// StreamChat 创建会话流式返回
func (c *Client) StreamChat(ctx context.Context, request *ChatRequest) (chan StreamChat, error) {
	oauth, err := c.GetOAuth()
	if err != nil {
		return nil, err
	}

	req := &coze.CreateChatsReq{
		ConversationID: request.ConversationID,
		BotID:          request.BotID,
		UserID:         request.UserID,
		Messages: []*coze.Message{
			coze.BuildUserQuestionText(request.Content, nil),
		},
	}

	reader, err := coze.NewCozeAPI(coze.NewJWTAuth(oauth, nil), coze.WithBaseURL(coze.CnBaseURL), coze.WithHttpClient(&http.Client{Timeout: time.Minute})).Chat.Stream(ctx, req)
	if err != nil {
		return nil, err
	}

	ch := make(chan StreamChat, 10)

	go func() {
		defer func() {
			_ = reader.Close()
			close(ch)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				event, err := reader.Recv()
				if errors.Is(err, io.EOF) {
					return
				}
				if err != nil {
					ch <- StreamChat{Err: err}
					return
				}
				if event.Event == coze.ChatEventConversationMessageDelta && len(event.Message.Content) > 0 {
					ch <- StreamChat{Content: event.Message.Content}
				}
			}
		}

	}()

	return ch, nil
}

// Retrieve 查看对话详情
func (c *Client) Retrieve(ctx context.Context, request *RetrieveRequest) (*Chat, error) {
	oauth, err := c.GetOAuth()
	if err != nil {
		return nil, err
	}
	req := &coze.RetrieveChatsReq{
		ConversationID: request.ConversationId,
		ChatID:         request.ChatID,
	}
	resp, err := coze.NewCozeAPI(coze.NewJWTAuth(oauth, nil), coze.WithBaseURL(coze.CnBaseURL)).Chat.Retrieve(ctx, req)
	if err != nil {
		return nil, err
	}
	return &Chat{ID: resp.ID, ConversationID: resp.ConversationID, BotID: resp.BotID, CreatedAt: resp.CreatedAt, CompletedAt: resp.CompletedAt, FailedAt: resp.FailedAt, MetaData: resp.MetaData, LastError: resp.LastError.Msg, Status: string(resp.Status)}, nil
}

// List 查看对话列表
func (c *Client) List(ctx context.Context, request *ListRequest) (*ListResponse, error) {
	oauth, err := c.GetOAuth()
	if err != nil {
		return nil, err
	}

	req := &coze.ListChatsMessagesReq{
		ConversationID: request.ConversationId,
		ChatID:         request.ChatID,
	}
	resp, err := coze.NewCozeAPI(coze.NewJWTAuth(oauth, nil), coze.WithBaseURL(coze.CnBaseURL)).Chat.Messages.List(ctx, req)
	if err != nil {
		return nil, err
	}

	items := make([]*Message, 0, len(resp.Messages))
	for _, item := range resp.Messages {
		items = append(items, &Message{
			Role: string(item.Role), Type: string(item.Type), Content: item.Content, ReasoningContent: item.ReasoningContent, ContentType: string(item.ContentType), MetaData: item.MetaData, ID: item.ID, ConversationID: item.ConversationID, SectionID: item.SectionID, BotID: item.BotID, ChatID: item.ChatID, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
		})
	}
	return &ListResponse{Items: items}, nil
}
