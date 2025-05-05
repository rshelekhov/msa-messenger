package subscriber

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

func (u *Usecase) GetFriends(ctx context.Context, pageToken *string) (friends []entity.Friend, nextPageToken string, err error)
func (u *Usecase) GetFriend(ctx context.Context, friendID string) (friend entity.Friend, err error)
func (u *Usecase) RemoveFriend(ctx context.Context, friendID string) (err error)
func (u *Usecase) InviteFriend(ctx context.Context, invitation entity.FriendInvite) (invitationID string, err error)
func (u *Usecase) GetFriendInvites(ctx context.Context, direction, status, pageToken *string) (invites []entity.FriendInvite, nextPageToken string, err error)
func (u *Usecase) AcceptFriendInvite(ctx context.Context, inviteID string) (err error)
func (u *Usecase) DeclineFriendInvite(ctx context.Context, inviteID string) (err error)
