package client

import (
	"context"
	"log/slog"

	"github.com/rshelekhov/msa-messenger/gateway/internal/config"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/entity"
	authv1 "github.com/rshelekhov/sso-protos/gen/go/api/auth/v1"
	commonv1 "github.com/rshelekhov/sso-protos/gen/go/api/common/v1"
	"github.com/rshelekhov/sso/pkg/grpcerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type AuthClient struct {
	log                *slog.Logger
	GRPCClient         authv1.AuthServiceClient
	VerificationURL    string
	ConfirmPasswordURL string
}

func NewAuthClient(log *slog.Logger, conn *grpc.ClientConn, ssoService config.SSOService) *AuthClient {
	return &AuthClient{
		log:                log,
		GRPCClient:         authv1.NewAuthServiceClient(conn),
		VerificationURL:    ssoService.VerificationURL,
		ConfirmPasswordURL: ssoService.ConfirmPasswordURL,
	}
}

func (a *AuthClient) RegisterUser(ctx context.Context, user entity.UserCredentials, device entity.UserDevice) (userID string, tokens entity.UserTokens, err error) {
	const op = "grpc.client.auth.RegisterUser"

	req := &authv1.RegisterUserRequest{
		Email:              user.Email,
		Password:           user.Password,
		Name:               user.Name,
		VerificationUrl:    a.VerificationURL,
		ConfirmPasswordUrl: a.ConfirmPasswordURL,
		UserDeviceData: &authv1.UserDeviceData{
			UserAgent: device.UserAgent,
			Ip:        device.IP,
		},
	}

	resp, err := a.GRPCClient.RegisterUser(ctx, req)
	if err != nil {
		return "", entity.UserTokens{}, a.mapSSOError(err)
	}

	userID = resp.GetUserId()
	if userID == "" {
		a.log.Error("invalid response from SSO service",
			slog.String("op", op),
			slog.String("issue", "user_id is empty"),
		)
		return "", entity.UserTokens{}, domain.ErrInvalidResponse
	}

	tokenData := resp.GetTokenData()
	if tokenData == nil {
		a.log.Error("invalid response from SSO service",
			slog.String("op", op),
			slog.String("issue", "token data is nil"),
		)
		return "", entity.UserTokens{}, domain.ErrInvalidResponse
	}

	return userID, entity.UserTokens{
		AccessToken:      tokenData.AccessToken,
		RefreshToken:     tokenData.RefreshToken,
		Domain:           tokenData.Domain,
		Path:             tokenData.Path,
		ExpiresAt:        tokenData.ExpiresAt.AsTime(),
		HttpOnly:         tokenData.HttpOnly,
		AdditionalFields: tokenData.AdditionalFields,
	}, nil
}

func (a *AuthClient) Login(ctx context.Context, user entity.UserCredentials, device entity.UserDevice) (userID string, tokens entity.UserTokens, err error) {
	const op = "grpc.client.auth.Login"

	req := &authv1.LoginRequest{
		Email:    user.Email,
		Password: user.Password,
		UserDeviceData: &authv1.UserDeviceData{
			UserAgent: device.UserAgent,
			Ip:        device.IP,
		},
	}

	resp, err := a.GRPCClient.Login(ctx, req)
	if err != nil {
		return "", entity.UserTokens{}, a.mapSSOError(err)
	}

	userID = resp.GetUserId()
	if userID == "" {
		a.log.Error("invalid response from SSO service",
			slog.String("op", op),
			slog.String("issue", "user_id is empty"),
		)
		return "", entity.UserTokens{}, domain.ErrInvalidResponse
	}

	tokenData := resp.GetTokenData()
	if tokenData == nil {
		a.log.Error("invalid response from SSO service",
			slog.String("op", op),
			slog.String("issue", "token data is nil"),
		)
		return "", entity.UserTokens{}, domain.ErrInvalidResponse
	}

	return userID, entity.UserTokens{
		AccessToken:      tokenData.AccessToken,
		RefreshToken:     tokenData.RefreshToken,
		Domain:           tokenData.Domain,
		Path:             tokenData.Path,
		ExpiresAt:        tokenData.ExpiresAt.AsTime(),
		HttpOnly:         tokenData.HttpOnly,
		AdditionalFields: tokenData.AdditionalFields,
	}, nil
}

func (a *AuthClient) Logout(ctx context.Context, device entity.UserDevice) (err error) {
	req := &authv1.LogoutRequest{
		UserDeviceData: &authv1.UserDeviceData{
			UserAgent: device.UserAgent,
			Ip:        device.IP,
		},
	}

	_, err = a.GRPCClient.Logout(ctx, req)
	if err != nil {
		return a.mapSSOError(err)
	}

	return nil
}
func (a *AuthClient) RefreshToken(ctx context.Context, refreshToken string, device entity.UserDevice) (tokens entity.UserTokens, err error) {
	req := &authv1.RefreshTokensRequest{
		RefreshToken: refreshToken,
		UserDeviceData: &authv1.UserDeviceData{
			UserAgent: device.UserAgent,
			Ip:        device.IP,
		},
	}

	resp, err := a.GRPCClient.RefreshTokens(ctx, req)
	if err != nil {
		return entity.UserTokens{}, a.mapSSOError(err)
	}

	return entity.UserTokens{
		AccessToken:      resp.GetTokenData().AccessToken,
		RefreshToken:     resp.GetTokenData().RefreshToken,
		Domain:           resp.GetTokenData().Domain,
		Path:             resp.GetTokenData().Path,
		ExpiresAt:        resp.GetTokenData().ExpiresAt.AsTime(),
		HttpOnly:         resp.GetTokenData().HttpOnly,
		AdditionalFields: resp.GetTokenData().AdditionalFields,
	}, nil
}

func (a *AuthClient) GetJWKS(ctx context.Context) (jwks entity.JWKS, err error) {
	req := &authv1.GetJWKSRequest{}

	resp, err := a.GRPCClient.GetJWKS(ctx, req)
	if err != nil {
		return entity.JWKS{}, a.mapSSOError(err)
	}

	keys := make([]entity.JWK, len(resp.GetJwks()))
	for i, jwk := range resp.GetJwks() {
		keys[i] = entity.JWK{
			Alg: jwk.GetAlg(),
			Kty: jwk.GetKty(),
			Use: jwk.GetUse(),
			Kid: jwk.GetKid(),
			N:   jwk.GetN(),
			E:   jwk.GetE(),
		}
	}

	return entity.JWKS{
		Keys: keys,
	}, nil
}

func (a *AuthClient) VerifyEmail(ctx context.Context, verificationToken string) (err error) {
	req := &authv1.VerifyEmailRequest{
		Token: verificationToken,
	}

	_, err = a.GRPCClient.VerifyEmail(ctx, req)
	if err != nil {
		return a.mapSSOError(err)
	}

	return nil
}

func (a *AuthClient) ResetPassword(ctx context.Context, email string) (err error) {
	req := &authv1.ResetPasswordRequest{
		ConfirmUrl: a.ConfirmPasswordURL,
		Email:      email,
	}

	_, err = a.GRPCClient.ResetPassword(ctx, req)
	if err != nil {
		return a.mapSSOError(err)
	}

	return nil
}

func (a *AuthClient) ChangePassword(ctx context.Context, passwordResetToken string, updatedPassword string) (err error) {
	req := &authv1.ChangePasswordRequest{
		Token:           passwordResetToken,
		UpdatedPassword: updatedPassword,
	}

	_, err = a.GRPCClient.ChangePassword(ctx, req)
	if err != nil {
		return a.mapSSOError(err)
	}

	return nil
}

func (a *AuthClient) mapSSOError(err error) error {
	extracted, extractErr := grpcerrors.ExtractError(err)
	if extractErr != nil {
		a.log.Error("failed to extract error",
			slog.String("op", "mapSSOError"),
			slog.String("error", err.Error()),
		)
		return ErrServiceFailure
	}

	// Check error code if details are available
	if extracted.HasDetails {
		switch extracted.ErrorCode {
		case commonv1.ErrorCode_ERROR_CODE_INVALID_CREDENTIALS:
			return domain.ErrInvalidCredentials

		case commonv1.ErrorCode_ERROR_CODE_SESSION_EXPIRED:
			return domain.ErrSessionExpired

		case commonv1.ErrorCode_ERROR_CODE_SESSION_NOT_FOUND:
			return domain.ErrSessionNotFound

		case commonv1.ErrorCode_ERROR_CODE_USER_ALREADY_EXISTS:
			return domain.ErrUserAlreadyExists

		case commonv1.ErrorCode_ERROR_CODE_EMAIL_ALREADY_TAKEN:
			return domain.ErrEmailAlreadyTaken

		case commonv1.ErrorCode_ERROR_CODE_USER_NOT_FOUND:
			return domain.ErrUserNotFound

		case commonv1.ErrorCode_ERROR_CODE_VALIDATION_ERROR,
			commonv1.ErrorCode_ERROR_CODE_PASSWORDS_DO_NOT_MATCH,
			commonv1.ErrorCode_ERROR_CODE_CLIENT_ID_NOT_ALLOWED:
			return domain.ErrInvalidRequest

		case commonv1.ErrorCode_ERROR_CODE_TOKEN_EXPIRED_EMAIL_RESENT:
			return domain.ErrTokenExpiredEmailResent

		case commonv1.ErrorCode_ERROR_CODE_FAILED_TO_SEND_VERIFICATION_EMAIL:
			return domain.ErrFailedToSendVerificationEmail

		case commonv1.ErrorCode_ERROR_CODE_FAILED_TO_SEND_RESET_PASSWORD_EMAIL:
			return domain.ErrFailedToSendResetPasswordEmail

		case commonv1.ErrorCode_ERROR_CODE_USER_DEVICE_NOT_FOUND:
			return domain.ErrUserDeviceNotRegisteredOrAlreadyLoggedOut
		}

	}

	// Fallback to gRPC code if no specific error code match
	switch extracted.GRPCCode {
	case codes.Unauthenticated:
		return domain.ErrInvalidCredentials

	case codes.NotFound:
		return domain.ErrUserNotFound

	case codes.AlreadyExists:
		return domain.ErrUserAlreadyExists

	case codes.InvalidArgument:
		return domain.ErrInvalidRequest

	case codes.Internal, codes.Unavailable:
		return ErrServiceFailure

	default:
		a.log.Error("unknown error",
			slog.String("op", "mapSSOError"),
			slog.String("error", err.Error()),
		)
		return ErrServiceFailure
	}
}
