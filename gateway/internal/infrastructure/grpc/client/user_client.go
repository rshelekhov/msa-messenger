package client

import (
	"context"
	"log/slog"

	"github.com/rshelekhov/msa-messenger/gateway/internal/domain"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/entity"
	userv1 "github.com/rshelekhov/sso-protos/gen/go/api/user/v1"
	"google.golang.org/grpc"
)

type User struct {
	log        *slog.Logger
	GRPCClient userv1.UserServiceClient
}

func NewUserClient(log *slog.Logger, conn *grpc.ClientConn) *User {
	return &User{
		log:        log,
		GRPCClient: userv1.NewUserServiceClient(conn),
	}
}
func (u *User) GetUser(ctx context.Context) (user entity.User, err error) {
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

func (u *User) UpdateUser(ctx context.Context, user entity.User) (userID string, err error) {
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

func (u *User) DeleteUser(ctx context.Context) (err error) {
	req := &userv1.DeleteUserRequest{}

	_, err = u.GRPCClient.DeleteUser(ctx, req)
	if err != nil {
		return u.mapUserError(err)
	}

	return nil
}

func (u *User) SearchUsers(ctx context.Context, query string) (users []entity.User, err error) {
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

func (u *User) mapUserError(err error) error {
	// TODO: implement this
	return nil
}
