package auth

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

func (u *Usecase) Register(ctx context.Context, user entity.UserCredentials, device entity.UserDevice) (userID string, tokens entity.UserTokens, err error)
func (u *Usecase) Login(ctx context.Context, user entity.UserCredentials, device entity.UserDevice) (tokens entity.UserTokens, err error)
func (u *Usecase) Logout(ctx context.Context, device entity.UserDevice) (err error)
func (u *Usecase) RefreshToken(ctx context.Context, refreshToken string, device entity.UserDevice) (tokens entity.UserTokens, err error)
func (u *Usecase) GetJWKS(ctx context.Context) (jwks entity.JWKS, err error)
