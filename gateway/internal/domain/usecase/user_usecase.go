package usecase

import (
	"context"

	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/entity"
)

type UserClient interface {
	GetUser(ctx context.Context) (user entity.User, err error)
	UpdateUser(ctx context.Context, user entity.User) (userID string, err error)
	DeleteUser(ctx context.Context) (err error)
	SearchUsers(ctx context.Context, query string) (users []entity.User, err error)
}

type UserUsecase struct {
	userClient UserClient
}

func NewUserUsecase(userClient UserClient) *UserUsecase {
	return &UserUsecase{userClient: userClient}
}

func (u *UserUsecase) GetOwnProfile(ctx context.Context) (user entity.User, err error) {
	return u.userClient.GetUser(ctx)
}

func (u *UserUsecase) UpdateOwnProfile(ctx context.Context, user entity.User) (userID string, err error) {
	return u.userClient.UpdateUser(ctx, user)
}

func (u *UserUsecase) DeleteOwnProfile(ctx context.Context) (err error) {
	return u.userClient.DeleteUser(ctx)
}

func (u *UserUsecase) SearchUsers(ctx context.Context, query string) (users []entity.User, err error) {
	return u.userClient.SearchUsers(ctx, query)
}
