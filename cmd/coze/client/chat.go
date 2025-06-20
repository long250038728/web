package client

import (
	"context"
	"errors"
	"github.com/coze-dev/coze-go"
	"io"
	"net/http"
	"time"
)

// Chat 创建会话
func (c *Client) Chat(ctx context.Context, request *ChatRequest) (*ListResponse, error) {
	req := &coze.CreateChatsReq{
		ConversationID: request.ConversationID,
		BotID:          request.BotID,
		UserID:         request.UserID,
		Messages: []*coze.Message{
			coze.BuildUserQuestionText(request.Content, nil),
		},
	}

	timeout := int(time.Second) * 20
	resp, err := c.getApi().Chat.CreateAndPoll(ctx, req, &timeout)
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
	req := &coze.CreateChatsReq{
		ConversationID: request.ConversationID,
		BotID:          request.BotID,
		UserID:         request.UserID,
		Messages: []*coze.Message{
			coze.BuildUserQuestionText(request.Content, nil),
		},
	}
	reader, err := c.getApi(coze.WithHttpClient(&http.Client{Timeout: time.Minute})).Chat.Stream(ctx, req)
	if err != nil {
		return nil, err
	}

	ch := make(chan StreamChat, 10)

	go func() {
		defer func() {
			_ = reader.Close()
			close(ch)
		}()

		isThinkStop := false

		for {
			select {
			case <-ctx.Done():
				return
			default:
				event, err := reader.Recv()
				if errors.Is(err, io.EOF) {
					ch <- StreamChat{Content: "\n </body>"}
					return
				}
				if err != nil {
					ch <- StreamChat{Err: err}
					return
				}
				if event.Event == coze.ChatEventConversationChatInProgress {
					ch <- StreamChat{Content: "<think>"}
				}

				if event.Event == coze.ChatEventConversationMessageDelta && len(event.Message.ReasoningContent) > 0 {
					ch <- StreamChat{Content: event.Message.ReasoningContent}
				}

				if event.Event == coze.ChatEventConversationMessageDelta && len(event.Message.Content) > 0 {
					if isThinkStop == false {
						isThinkStop = true
						ch <- StreamChat{Content: "</think> \n"}
						ch <- StreamChat{Content: "<body>"}
					}

					ch <- StreamChat{Content: event.Message.Content}
				}
			}
		}

	}()

	return ch, nil
}

// Retrieve 查看对话详情
func (c *Client) Retrieve(ctx context.Context, request *RetrieveRequest) (*ChatItem, error) {
	req := &coze.RetrieveChatsReq{
		ConversationID: request.ConversationID,
		ChatID:         request.ChatID,
	}
	resp, err := c.getApi().Chat.Retrieve(ctx, req)
	if err != nil {
		return nil, err
	}
	chat := &ChatItem{
		ID:             resp.ID,
		ConversationID: resp.ConversationID,
		BotID:          resp.BotID,
		CreatedAt:      resp.CreatedAt,
		CompletedAt:    resp.CompletedAt,
		FailedAt:       resp.FailedAt,
		MetaData:       resp.MetaData,
		Status:         string(resp.Status)}

	if resp.LastError != nil {
		chat.LastError = resp.LastError.Msg
	}

	return chat, nil
}

// List 查看对话列表
func (c *Client) List(ctx context.Context, request *ListRequest) (*ListResponse, error) {
	req := &coze.ListChatsMessagesReq{
		ConversationID: request.ConversationID,
		ChatID:         request.ChatID,
	}
	resp, err := c.getApi().Chat.Messages.List(ctx, req)
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
