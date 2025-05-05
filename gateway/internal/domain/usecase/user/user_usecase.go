package user

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

// TODO: add here usecase for updating user profile — need to have client for auth and user services
func (u *Usecase) GetOwnProfile(ctx context.Context) (user entity.User, err error)
func (u *Usecase) UpdateOwnProfile(ctx context.Context, user entity.User) (userID string, err error)
func (u *Usecase) DeleteOwnProfile(ctx context.Context) (err error)
func (u *Usecase) SearchUsers(ctx context.Context, query string) (users []entity.User, err error)
