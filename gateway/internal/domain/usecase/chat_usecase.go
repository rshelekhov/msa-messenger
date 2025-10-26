package usecase

import (
	"context"

	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/entity"
)

type ChatClient interface {
	CreateChat(ctx context.Context, friendID string) (chatID string, err error)
	GetChat(ctx context.Context, chatID string, pageToken *string) (chat entity.Chat, nextPageToken string, err error)
	DeleteChat(ctx context.Context, chatID string) (err error)
	ListChats(ctx context.Context, pageToken *string) (chats []entity.Chat, nextPageToken string, err error)
	SendMessage(ctx context.Context, chatID string, content string) (message entity.Message, err error)
}

type ChatUsecase struct {
	chatClient ChatClient
}

func NewChatUsecase(chatClient ChatClient) *ChatUsecase {
	return &ChatUsecase{chatClient: chatClient}
}

func (u *ChatUsecase) GetChats(ctx context.Context, pageToken *string) (chats []entity.Chat, nextPageToken string, err error) {
	return u.chatClient.ListChats(ctx, pageToken)
}

func (u *ChatUsecase) GetChat(ctx context.Context, chatID string, pageToken *string) (chat entity.Chat, nextPageToken string, err error) {
	return u.chatClient.GetChat(ctx, chatID, pageToken)
}

func (u *ChatUsecase) SendMessage(ctx context.Context, chatID string, content string) (message entity.Message, err error) {
	return u.chatClient.SendMessage(ctx, chatID, content)
}

func (u *ChatUsecase) DeleteChat(ctx context.Context, chatID string) (err error) {
	return u.chatClient.DeleteChat(ctx, chatID)
}
