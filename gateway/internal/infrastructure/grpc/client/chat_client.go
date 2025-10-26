package client

import (
	"context"
	"log/slog"

	chatv1 "github.com/rshelekhov/msa-messenger-protos/gen/go/api/chat/v1"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/entity"
	"google.golang.org/grpc"
)

type ChatClient struct {
	log        *slog.Logger
	GRPCClient chatv1.ChatServiceClient
}

func NewChatClient(log *slog.Logger, conn *grpc.ClientConn) *ChatClient {
	return &ChatClient{
		log:        log,
		GRPCClient: chatv1.NewChatServiceClient(conn),
	}
}

func (c *ChatClient) CreateChat(ctx context.Context, friendID string) (chatID string, err error) {
	const op = "grpc.client.chat.CreateChat"

	req := &chatv1.CreateChatRequest{
		FriendId: friendID,
	}

	resp, err := c.GRPCClient.CreateChat(ctx, req)
	if err != nil {
		return "", c.mapChatError(err)
	}

	chatID = resp.GetChat().GetId()
	if chatID == "" {
		c.log.Error("invalid response from chat service",
			slog.String("op", op),
			slog.String("issue", "chat id is empty"),
		)
		return "", domain.ErrInvalidResponse
	}

	return chatID, nil
}

func (c *ChatClient) GetChat(ctx context.Context, chatID string, pageToken *string) (chat entity.Chat, nextPageToken string, err error) {
	const op = "grpc.client.chat.GetChat"

	req := &chatv1.GetChatRequest{
		ChatId: chatID,
	}
	if pageToken != nil {
		req.PageToken = *pageToken
	}

	resp, err := c.GRPCClient.GetChat(ctx, req)
	if err != nil {
		return entity.Chat{}, "", c.mapChatError(err)
	}

	chatData := resp.GetChat()
	if chatData == nil {
		c.log.Error("invalid response from chat service",
			slog.String("op", op),
			slog.String("issue", "chat data is nil"),
		)
		return entity.Chat{}, "", domain.ErrInvalidResponse
	}

	chat = entity.Chat{
		ID:        chatData.GetId(),
		UserID:    chatData.GetUserId(),
		FriendID:  chatData.GetFriendId(),
		CreatedAt: chatData.GetCreatedAt().AsTime(),
		Messages:  mapMessages(resp.GetMessages()),
	}

	return chat, resp.GetNextPageToken(), nil
}

func mapMessages(protoMessages []*chatv1.Message) []entity.Message {
	if len(protoMessages) == 0 {
		return []entity.Message{}
	}

	messages := make([]entity.Message, len(protoMessages))
	for i, msg := range protoMessages {
		messages[i] = entity.Message{
			ID:        msg.GetId(),
			SenderID:  msg.GetSenderId(),
			Content:   msg.GetContent(),
			CreatedAt: msg.GetCreatedAt().AsTime(),
		}
	}
	return messages
}

func (c *ChatClient) DeleteChat(ctx context.Context, chatID string) (err error) {
	req := &chatv1.DeleteChatRequest{
		ChatId: chatID,
	}

	_, err = c.GRPCClient.DeleteChat(ctx, req)
	if err != nil {
		return c.mapChatError(err)
	}

	return nil
}

func (c *ChatClient) ListChats(ctx context.Context, pageToken *string) (chats []entity.Chat, nextPageToken string, err error) {
	req := &chatv1.ListChatsRequest{}
	if pageToken != nil {
		req.PageToken = *pageToken
	}

	resp, err := c.GRPCClient.ListChats(ctx, req)
	if err != nil {
		return nil, "", c.mapChatError(err)
	}

	chatsData := resp.GetChats()
	chats = make([]entity.Chat, len(chatsData))
	for i, chatData := range chatsData {
		chats[i] = entity.Chat{
			ID:        chatData.GetId(),
			UserID:    chatData.GetUserId(),
			FriendID:  chatData.GetFriendId(),
			CreatedAt: chatData.GetCreatedAt().AsTime(),
			Messages:  []entity.Message{}, // Messages will be loaded separately if needed
		}
	}

	return chats, resp.GetNextPageToken(), nil
}

func (c *ChatClient) SendMessage(ctx context.Context, chatID string, content string) (message entity.Message, err error) {
	const op = "grpc.client.chat.SendMessage"

	req := &chatv1.SendMessageRequest{
		ChatId:  chatID,
		Content: content,
	}

	resp, err := c.GRPCClient.SendMessage(ctx, req)
	if err != nil {
		return entity.Message{}, c.mapChatError(err)
	}

	msgData := resp.GetMessage()
	if msgData == nil {
		c.log.Error("invalid response from chat service",
			slog.String("op", op),
			slog.String("issue", "message data is nil"),
		)
		return entity.Message{}, domain.ErrInvalidResponse
	}

	message = entity.Message{
		ID:        msgData.GetId(),
		SenderID:  msgData.GetSenderId(),
		Content:   msgData.GetContent(),
		CreatedAt: msgData.GetCreatedAt().AsTime(),
	}

	return message, nil
}

func (c *ChatClient) mapChatError(err error) error {
	// TODO: implement this
	return nil
}
