package client

import (
	"context"
	"log/slog"

	"github.com/rshelekhov/msa-messenger/gateway/internal/domain"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/entity"
	commonv1 "github.com/rshelekhov/sso-protos/gen/go/api/common/v1"
	userv1 "github.com/rshelekhov/sso-protos/gen/go/api/user/v1"
	"github.com/rshelekhov/sso/pkg/grpcerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type UserClient struct {
	log        *slog.Logger
	GRPCClient userv1.UserServiceClient
}

func NewUserClient(log *slog.Logger, conn *grpc.ClientConn) *UserClient {
	return &UserClient{
		log:        log,
		GRPCClient: userv1.NewUserServiceClient(conn),
	}
}
func (u *UserClient) GetUser(ctx context.Context) (user entity.User, err error) {
	const op = "grpc.client.user.GetUser"

	req := &userv1.GetUserRequest{}

	resp, err := u.GRPCClient.GetUser(ctx, req)
	if err != nil {
		return entity.User{}, u.mapUserError(err)
	}

	userData := resp.GetUser()
	if userData == nil {
		u.log.Error("invalid response from user service",
			slog.String("op", op),
			slog.String("issue", "user data is nil"),
		)
		return entity.User{}, domain.ErrInvalidResponse
	}

	return entity.User{
		ID:    userData.GetId(),
		Email: userData.GetEmail(),
		Name:  userData.GetName(),
	}, nil
}

func (u *UserClient) UpdateUser(ctx context.Context, user entity.User) (userID string, err error) {
	req := &userv1.UpdateUserRequest{
		Email:           user.Email,
		Name:            user.Name,
		CurrentPassword: user.CurrentPassword,
		UpdatedPassword: user.UpdatedPassword,
	}

	_, err = u.GRPCClient.UpdateUser(ctx, req)
	if err != nil {
		return "", u.mapUserError(err)
	}

	return user.ID, nil
}

func (u *UserClient) DeleteUser(ctx context.Context) (err error) {
	req := &userv1.DeleteUserRequest{}

	_, err = u.GRPCClient.DeleteUser(ctx, req)
	if err != nil {
		return u.mapUserError(err)
	}

	return nil
}

func (u *UserClient) SearchUsers(ctx context.Context, query string) (users []entity.User, err error) {
	req := &userv1.SearchUsersRequest{
		Query: query,
	}

	resp, err := u.GRPCClient.SearchUsers(ctx, req)
	if err != nil {
		return nil, u.mapUserError(err)
	}

	usersData := resp.GetUsers()
	if usersData == nil {
		return nil, domain.ErrNoUsersFound
	}

	users = make([]entity.User, len(usersData))
	for i, userData := range usersData {
		users[i] = entity.User{
			ID:    userData.GetId(),
			Email: userData.GetEmail(),
		}
	}

	return users, nil
}

func (u *UserClient) mapUserError(err error) error {
	extracted, extractErr := grpcerrors.ExtractError(err)
	if extractErr != nil {
		u.log.Error("failed to extract error",
			slog.String("op", "mapSSOError"),
			slog.String("error", err.Error()),
		)
		return ErrServiceFailure
	}

	// Check error code if details are available
	if extracted.HasDetails {
		switch extracted.ErrorCode {
		case commonv1.ErrorCode_ERROR_CODE_USER_NOT_FOUND:
			return domain.ErrUserNotFound

		case commonv1.ErrorCode_ERROR_CODE_VALIDATION_ERROR:
			return domain.ErrInvalidRequest

		case commonv1.ErrorCode_ERROR_CODE_CURRENT_PASSWORD_REQUIRED:
			return domain.ErrCurrentPasswordRequired

		case commonv1.ErrorCode_ERROR_CODE_NO_EMAIL_CHANGES_DETECTED:
			return domain.ErrNoEmailChangesDetected

		case commonv1.ErrorCode_ERROR_CODE_NO_PASSWORD_CHANGES_DETECTED:
			return domain.ErrNoPasswordChangesDetected

		case commonv1.ErrorCode_ERROR_CODE_NO_NAME_CHANGES_DETECTED:
			return domain.ErrNoNameChangesDetected

		case commonv1.ErrorCode_ERROR_CODE_PASSWORDS_DO_NOT_MATCH:
			return domain.ErrPasswordsDoNotMatch

		case commonv1.ErrorCode_ERROR_CODE_INVALID_CREDENTIALS:
			return domain.ErrCurrentPasswordIsIncorrect

		case commonv1.ErrorCode_ERROR_CODE_EMAIL_ALREADY_TAKEN:
			return domain.ErrEmailAlreadyTaken
		}
	}

	// Fallback to gRPC code if no specific error code match
	switch extracted.GRPCCode {
	case codes.NotFound:
		return domain.ErrUserNotFound

	case codes.InvalidArgument:
		return domain.ErrInvalidRequest

	case codes.Internal, codes.Unavailable:
		return ErrServiceFailure

	default:
		u.log.Error("unknown error",
			slog.String("op", "mapUserError"),
			slog.String("error", err.Error()),
		)
		return ErrServiceFailure
	}
}
