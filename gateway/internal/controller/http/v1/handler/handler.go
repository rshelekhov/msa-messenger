package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/entity"
)

type Handler struct {
	log               *slog.Logger
	validate          *validator.Validate
	tokenSender       *tokenSender
	tokenManager      TokenManager
	authUsecase       AuthUsecase
	userUsecase       UserUsecase
	subscriberUsecase SubscriberUsecase
	chatUsecase       ChatUsecase
}

type (
	TokenManager interface {
		ExtractRefreshTokenFromCookies(r *http.Request) (string, error)
	}

	AuthUsecase interface {
		Register(ctx context.Context, user entity.UserCredentials, device entity.UserDevice) (userID string, tokens entity.UserTokens, err error)
		Login(ctx context.Context, user entity.UserCredentials, device entity.UserDevice) (userID string, tokens entity.UserTokens, err error)
		Logout(ctx context.Context, device entity.UserDevice) (err error)
		RefreshToken(ctx context.Context, refreshToken string, device entity.UserDevice) (tokens entity.UserTokens, err error)
		GetJWKS(ctx context.Context) (jwks entity.JWKS, err error)
		VerifyEmail(ctx context.Context, verificationToken string) (err error)
		ResetPassword(ctx context.Context, email string) (err error)
		ChangePassword(ctx context.Context, passwordResetToken string, updatedPassword string) (err error)
	}

	UserUsecase interface {
		GetOwnProfile(ctx context.Context) (user entity.User, err error)
		UpdateOwnProfile(ctx context.Context, user entity.User) (userID string, err error)
		DeleteOwnProfile(ctx context.Context) (err error)
		SearchUsers(ctx context.Context, query string) (users []entity.User, err error)
	}

	SubscriberUsecase interface {
		GetFriends(ctx context.Context, pageToken *string) (friends []entity.Friend, nextPageToken string, err error)
		GetFriend(ctx context.Context, friendID string) (friend entity.Friend, err error)
		RemoveFriend(ctx context.Context, friendID string) (err error)
		InviteFriend(ctx context.Context, invitation entity.FriendInvite) (invitationID string, err error)
		GetFriendInvites(ctx context.Context, direction, status, pageToken *string) (invites []entity.FriendInvite, nextPageToken string, err error)
		AcceptFriendInvite(ctx context.Context, inviteID string) (err error)
		DeclineFriendInvite(ctx context.Context, inviteID string) (err error)
	}

	ChatUsecase interface {
		GetChats(ctx context.Context, pageToken *string) (chats []entity.Chat, nextPageToken string, err error)
		GetChat(ctx context.Context, chatID string, pageToken *string) (chat entity.Chat, nextPageToken string, err error)
		SendMessage(ctx context.Context, chatID string, content string) (message entity.Message, err error)
		DeleteChat(ctx context.Context, chatID string) (err error)
	}
)

func New(
	log *slog.Logger,
	validate *validator.Validate,
	tokenManager TokenManager,
	authUsecase AuthUsecase,
	userUsecase UserUsecase,
	subscriberUsecase SubscriberUsecase,
	chatUsecase ChatUsecase,
) *Handler {
	return &Handler{
		log:               log,
		validate:          validate,
		tokenSender:       NewTokenSender(),
		tokenManager:      tokenManager,
		authUsecase:       authUsecase,
		userUsecase:       userUsecase,
		subscriberUsecase: subscriberUsecase,
		chatUsecase:       chatUsecase,
	}
}

func (h *Handler) logWithReqID(ctx context.Context, op string) *slog.Logger {
	reqID := middleware.GetReqID(ctx)
	return h.log.With(
		slog.String("req_id", reqID),
		slog.String("op", op),
	)
}

func (h *Handler) decodeAndValidateRequest(w http.ResponseWriter, r *http.Request, log *slog.Logger, request any) error {
	if err := render.Decode(r, request); err != nil {
		log.Warn("failed to decode request",
			slog.String("error", err.Error()),
		)
		handleBadRequestError(w, r, ErrInvalidRequest)
		return err
	}

	if err := h.validate.Struct(request); err != nil {
		handleValidationErrors(w, r, log, err)
		return err
	}

	return nil
}
