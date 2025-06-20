package client

import (
	"context"
	"github.com/coze-dev/coze-go"
)

type CozeClientInterface interface {
	GetAccessToken(ctx context.Context) (string, error)
	ConversationsInterface
	ConversationsMessageInterface
	ChatInterface
}

type ConversationsInterface interface {
	ConversationsCreate(ctx context.Context, request *ConversationsCreateRequest) (*Conversation, error)
	ConversationsList(ctx context.Context, request *ConversationsListRequest) (*ConversationsListResponse, error)
	ConversationsRetrieve(ctx context.Context, request *ConversationsRetrieveRequest) (*Conversation, error)
	ConversationsClear(ctx context.Context, request *ConversationsClearRequest) (string, error)
}
type ConversationsMessageInterface interface {
	ConversationsMessageCreate(ctx context.Context, request *ConversationsMessageCreateRequest) (*Message, error)
	ConversationsMessageList(ctx context.Context, request *ConversationsMessageListRequest) (*ConversationsMessageListResponse, error)
	ConversationsMessageRetrieve(ctx context.Context, request *ConversationsMessageRetrieveRequest) (*coze.RetrieveConversationsMessagesResp, error)
}
type ChatInterface interface {
	Chat(ctx context.Context, request *ChatRequest) (*ListResponse, error)
	StreamChat(ctx context.Context, request *ChatRequest) (chan StreamChat, error)
	Retrieve(ctx context.Context, request *RetrieveRequest) (*ChatItem, error)
	List(ctx context.Context, request *ListRequest) (*ListResponse, error)
}

// ====================================================================

type Conversation struct {
	ID            string            `json:"id"`
	CreatedAt     int               `json:"created_at"`
	MetaData      map[string]string `json:"meta_data"`
	LastSectionID string            `json:"last_section_id"`
}
type ChatItem struct {
	ID             string            `json:"id"`
	ConversationID string            `json:"conversation_id"`
	BotID          string            `json:"bot_id"`
	CreatedAt      int               `json:"created_at"`
	CompletedAt    int               `json:"completed_at"`
	FailedAt       int               `json:"failed_at"`
	MetaData       map[string]string `json:"meta_data"`
	LastError      string            `json:"last_error"`
	Status         string            `json:"status"`
}
type Message struct {
	Role             string            `json:"role"`
	Type             string            `json:"type"`
	Content          string            `json:"content"`
	ReasoningContent string            `json:"reasoning_content"`
	ContentType      string            `json:"content_type"`
	MetaData         map[string]string `json:"meta_data"`
	ID               string            `json:"id"`
	ConversationID   string            `json:"conversation_id"`
	SectionID        string            `json:"section_id"`
	BotID            string            `json:"bot_id"`
	ChatID           string            `json:"chat_id"`
	CreatedAt        int64             `json:"created_at"`
	UpdatedAt        int64             `json:"updated_at"`
}
type StreamChat struct {
	Content string
	Err     error
}

//====================================================================

type ConversationsListRequest struct {
	BotID    string `json:"bot_id"`
	PageNum  int    `json:"page_num"`
	PageSize int    `json:"page_size"`
}
type ConversationsRetrieveRequest struct {
	ConversationID string `json:"conversation_id"`
}
type ConversationsClearRequest struct {
	ConversationID string `json:"conversation_id"`
}
type ConversationsListResponse struct {
	Total int
	Items []*Conversation
}
type ConversationsCreateRequest struct {
	BotID    string            `json:"bot_id"`
	Content  string            `json:"content"`
	MetaData map[string]string `json:"meta_data"`
}

//====================================================================

type ConversationsMessageCreateRequest struct {
	ConversationID string
	Content        string
}
type ConversationsMessageListRequest struct {
	ConversationID string
}
type ConversationsMessageRetrieveRequest struct {
	ConversationID string
	MessageID      string
}
type ConversationsMessageListResponse struct {
	Items []*Message
}

// ====================================================================

type ChatRequest struct {
	ConversationID string `json:"conversation_id"`
	BotID          string `json:"bot_id"`
	UserID         string `json:"user_id"`

	Content  string            `json:"content"`
	MetaData map[string]string `json:"meta_data"`
}
type RetrieveRequest struct {
	ConversationID string `json:"conversation_id"`
	ChatID         string `json:"chat_id"`
}
type ListRequest struct {
	ConversationID string `json:"conversation_id"`
	ChatID         string `json:"chat_id"`
}
type ListResponse struct {
	Items []*Message
}
type ChatResponse struct {
	Items []*Message
}
