package chat

import (
	"context"
	"log/slog"

	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/entity"
)

type Usecase struct {
	log *slog.Logger
}

func NewUsecase(log *slog.Logger) *Usecase {
	return &Usecase{log: log}
}

func (u *Usecase) GetChats(ctx context.Context, pageToken *string) (chats []entity.Chat, nextPageToken string, err error)
func (u *Usecase) GetChat(ctx context.Context, chatID string, pageToken *string) (chat entity.Chat, nextPageToken string, err error)
func (u *Usecase) SendMessage(ctx context.Context, chatID string, content string) (message entity.Message, err error)
func (u *Usecase) DeleteChat(ctx context.Context, chatID string) (err error)
