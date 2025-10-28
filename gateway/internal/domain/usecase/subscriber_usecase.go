package usecase

import (
	"context"

	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/entity"
)

type SubscriberClient interface {
	GetFriends(ctx context.Context, pageToken string) (friends []entity.Friend, nextPageToken string, err error)
	GetFriend(ctx context.Context, friendID string) (friend entity.Friend, err error)
	RemoveFriend(ctx context.Context, friendID string) (err error)
	InviteFriend(ctx context.Context, invitation entity.FriendInvite) (inviteID string, err error)
	GetFriendInvites(ctx context.Context, direction, status, pageToken *string) (invites []entity.FriendInvite, nextPageToken string, err error)
	AcceptFriendInvite(ctx context.Context, inviteID string) (err error)
	DeclineFriendInvite(ctx context.Context, inviteID string) (err error)
}

type SubscriberUsecase struct {
	subscriberClient SubscriberClient
}

func NewSubscriberUsecase(subscriberClient SubscriberClient) *SubscriberUsecase {
	return &SubscriberUsecase{subscriberClient: subscriberClient}
}

func (u *SubscriberUsecase) GetFriends(ctx context.Context, pageToken *string) (friends []entity.Friend, nextPageToken string, err error) {
	pageTokenStr := ""
	if pageToken != nil {
		pageTokenStr = *pageToken
	}

	return u.subscriberClient.GetFriends(ctx, pageTokenStr)
}

func (u *SubscriberUsecase) GetFriend(ctx context.Context, friendID string) (friend entity.Friend, err error) {
	return u.subscriberClient.GetFriend(ctx, friendID)
}

func (u *SubscriberUsecase) RemoveFriend(ctx context.Context, friendID string) (err error) {
	return u.subscriberClient.RemoveFriend(ctx, friendID)
}

func (u *SubscriberUsecase) InviteFriend(ctx context.Context, invitation entity.FriendInvite) (inviteID string, err error) {
	return u.subscriberClient.InviteFriend(ctx, invitation)
}

func (u *SubscriberUsecase) GetFriendInvites(ctx context.Context, direction, status, pageToken *string) (invites []entity.FriendInvite, nextPageToken string, err error) {
	return u.subscriberClient.GetFriendInvites(ctx, direction, status, pageToken)
}

func (u *SubscriberUsecase) AcceptFriendInvite(ctx context.Context, inviteID string) (err error) {
	return u.subscriberClient.AcceptFriendInvite(ctx, inviteID)
}

func (u *SubscriberUsecase) DeclineFriendInvite(ctx context.Context, inviteID string) (err error) {
	return u.subscriberClient.DeclineFriendInvite(ctx, inviteID)
}
