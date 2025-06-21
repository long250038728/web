package handle

import (
	"context"
	"github.com/long250038728/web/cmd/coze/client"
)

type Handle struct {
	cli client.CozeClientInterface
}

type ConversationsMessageListRequest struct {
	UserID int32 `json:"user_id"`
}

type StreamChatRequest struct {
	UserID         int32  `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	Content        string `json:"content"`
}

type ConversationsClearRequest struct {
	UserID         int32  `json:"user_id"`
	ConversationID string `json:"conversation_id"`
}

type ConversationsMessageListResponse struct {
	Items          []*ConversationsMessageItem `json:"items"`
	ConversationID string                      `json:"conversation_id"`
}
type ConversationsMessageItem struct {
	Role             string `json:"role"`
	Type             string `json:"type"`
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content"`
	ContentType      string `json:"content_type"`
	ID               string `json:"id"`
	ConversationID   string `json:"conversation_id"`
	SectionID        string `json:"section_id"`
	BotID            string `json:"bot_id"`
	ChatID           string `json:"chat_id"`
}

//====================================================================

var BotID = "7479292866154561548"

func NewHandle(cli client.CozeClientInterface) *Handle {
	return &Handle{cli: cli}
}

func (h *Handle) ConversationsMessageList(ctx context.Context, request *ConversationsMessageListRequest) (*ConversationsMessageListResponse, error) {
	//查询是否有ConversationID，没有新增一个
	//h.cli.ConversationsCreate(ctx, request)
	var ConversationID = "7517117256539881526"

	req := &client.ConversationsMessageListRequest{}
	req.ConversationID = ConversationID

	resp, err := h.cli.ConversationsMessageList(ctx, req)
	if err != nil {
		return nil, err
	}

	items := make([]*ConversationsMessageItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &ConversationsMessageItem{
			Role:             item.Role,
			Type:             item.Type,
			Content:          item.Content,
			ReasoningContent: item.ReasoningContent,
			ContentType:      item.ContentType,
		})
	}

	return &ConversationsMessageListResponse{Items: items, ConversationID: ConversationID}, nil
}

func (h *Handle) StreamChat(ctx context.Context, request *StreamChatRequest) (chan client.StreamChat, error) {

	var ConversationID = "7517117256539881526"

	req := &client.ChatRequest{
		ConversationID: ConversationID,
		BotID:          BotID,
		UserID:         "12345",
		Content:        "如何定制方案",
	}
	return h.cli.StreamChat(ctx, req)
}

//====================================================================
//
//func (h *Handle) ConversationsCreate(ctx context.Context, request *client.ConversationsCreateRequest) (*client.Conversation, error) {
//	return h.cli.ConversationsCreate(ctx, request)
//}

func (h *Handle) ConversationsClear(ctx context.Context, request *ConversationsClearRequest) (string, error) {
	req := &client.ConversationsClearRequest{
		ConversationID: request.ConversationID,
	}
	return h.cli.ConversationsClear(ctx, req)
}
