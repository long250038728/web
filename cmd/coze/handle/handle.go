package handle

import (
	"context"
	"github.com/coze-dev/coze-go"
	"github.com/long250038728/web/cmd/coze/client"
)

type Handle struct {
	cli *client.Client
}

//====================================================================

func NewHandle(cli *client.Client) client.CozeClientInterface {
	return &Handle{cli: cli}
}

func (h *Handle) GetAccessToken(ctx context.Context) (string, error) {
	return h.cli.GetAccessToken(ctx)
}

func (h *Handle) ConversationsCreate(ctx context.Context, request *client.ConversationsCreateRequest) (*client.Conversation, error) {
	return h.cli.ConversationsCreate(ctx, request)
}

func (h *Handle) ConversationsList(ctx context.Context, request *client.ConversationsListRequest) (*client.ConversationsListResponse, error) {
	return h.cli.ConversationsList(ctx, request)
}

func (h *Handle) ConversationsRetrieve(ctx context.Context, request *client.ConversationsRetrieveRequest) (*client.Conversation, error) {
	return h.cli.ConversationsRetrieve(ctx, request)
}

func (h *Handle) ConversationsClear(ctx context.Context, request *client.ConversationsClearRequest) (string, error) {
	return h.cli.ConversationsClear(ctx, request)
}

func (h *Handle) ConversationsMessageCreate(ctx context.Context, request *client.ConversationsMessageCreateRequest) (*client.Message, error) {
	return h.cli.ConversationsMessageCreate(ctx, request)
}

func (h *Handle) ConversationsMessageList(ctx context.Context, request *client.ConversationsMessageListRequest) (*client.ConversationsMessageListResponse, error) {
	return h.cli.ConversationsMessageList(ctx, request)
}

func (h *Handle) ConversationsMessageRetrieve(ctx context.Context, request *client.ConversationsMessageRetrieveRequest) (*coze.RetrieveConversationsMessagesResp, error) {
	return h.cli.ConversationsMessageRetrieve(ctx, request)
}

func (h *Handle) Chat(ctx context.Context, request *client.ChatRequest) (*client.ListResponse, error) {
	return h.cli.Chat(ctx, request)
}

func (h *Handle) StreamChat(ctx context.Context, request *client.ChatRequest) (chan client.StreamChat, error) {
	return h.cli.StreamChat(ctx, request)
}

func (h *Handle) Retrieve(ctx context.Context, request *client.RetrieveRequest) (*client.ChatItem, error) {
	return h.cli.Retrieve(ctx, request)
}

func (h *Handle) List(ctx context.Context, request *client.ListRequest) (*client.ListResponse, error) {
	return h.cli.List(ctx, request)
}
