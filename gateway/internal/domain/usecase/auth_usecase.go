package usecase

import (
	"context"

	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/entity"
)

type AuthClient interface {
	RegisterUser(ctx context.Context, user entity.UserCredentials, device entity.UserDevice) (userID string, tokens entity.UserTokens, err error)
	Login(ctx context.Context, user entity.UserCredentials, device entity.UserDevice) (userID string, tokens entity.UserTokens, err error)
	Logout(ctx context.Context, device entity.UserDevice) (err error)
	RefreshToken(ctx context.Context, refreshToken string, device entity.UserDevice) (tokens entity.UserTokens, err error)
	GetJWKS(ctx context.Context) (jwks entity.JWKS, err error)
	VerifyEmail(ctx context.Context, verificationToken string) (err error)
	ResetPassword(ctx context.Context, email string) (err error)
	ChangePassword(ctx context.Context, passwordResetToken string, updatedPassword string) (err error)
}

type AuthUsecase struct {
	authClient AuthClient
}

func NewAuthUsecase(authClient AuthClient) *AuthUsecase {
	return &AuthUsecase{authClient: authClient}
}

func (u *AuthUsecase) Register(
	ctx context.Context,
	user entity.UserCredentials,
	device entity.UserDevice,
) (
	userID string,
	tokens entity.UserTokens,
	err error,
) {
	return u.authClient.RegisterUser(ctx, user, device)
}

func (u *AuthUsecase) Login(ctx context.Context, user entity.UserCredentials, device entity.UserDevice) (userID string, tokens entity.UserTokens, err error) {
	return u.authClient.Login(ctx, user, device)
}

func (u *AuthUsecase) Logout(ctx context.Context, device entity.UserDevice) (err error) {
	return u.authClient.Logout(ctx, device)
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string, device entity.UserDevice) (tokens entity.UserTokens, err error) {
	return u.authClient.RefreshToken(ctx, refreshToken, device)
}

func (u *AuthUsecase) GetJWKS(ctx context.Context) (jwks entity.JWKS, err error) {
	return u.authClient.GetJWKS(ctx)
}

func (u *AuthUsecase) VerifyEmail(ctx context.Context, verificationToken string) (err error) {
	return u.authClient.VerifyEmail(ctx, verificationToken)
}

func (u *AuthUsecase) ResetPassword(ctx context.Context, email string) (err error) {
	return u.authClient.ResetPassword(ctx, email)
}

func (u *AuthUsecase) ChangePassword(ctx context.Context, passwordResetToken string, updatedPassword string) (err error) {
	return u.authClient.ChangePassword(ctx, passwordResetToken, updatedPassword)
}
