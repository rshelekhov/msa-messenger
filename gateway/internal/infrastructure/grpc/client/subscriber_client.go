package client

import (
	"context"
	"log/slog"

	subscriberv1 "github.com/rshelekhov/msa-messenger-protos/gen/go/api/subscriber/v1"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/entity"
	"google.golang.org/grpc"
)

type SubscriberClient struct {
	log        *slog.Logger
	GRPCClient subscriberv1.SubscriberServiceClient
}

func NewSubscriberClient(log *slog.Logger, conn *grpc.ClientConn) *SubscriberClient {
	return &SubscriberClient{
		log:        log,
		GRPCClient: subscriberv1.NewSubscriberServiceClient(conn),
	}
}

func (s *SubscriberClient) GetFriends(ctx context.Context, pageToken string) (friends []entity.Friend, nextPageToken string, err error) {
	req := &subscriberv1.GetFriendsRequest{
		PageToken: pageToken,
	}

	resp, err := s.GRPCClient.GetFriends(ctx, req)
	if err != nil {
		return nil, "", s.mapSubscriberError(err)
	}

	friendsData := resp.GetFriends()
	if friendsData == nil {
		return nil, "", domain.ErrNoFriends
	}

	friends = make([]entity.Friend, len(friendsData))
	for i, friendData := range friendsData {
		friends[i] = entity.Friend{
			ID:           friendData.GetUserId(),
			Name:         friendData.GetName(),
			FriendsSince: friendData.GetFriendsSince().AsTime(),
		}
	}
	nextPageToken = resp.GetNextPageToken()

	return friends, nextPageToken, nil
}

func (s *SubscriberClient) GetFriend(ctx context.Context, friendID string) (friend entity.Friend, err error) {
	const op = "grpc.client.subscriber.GetFriend"

	req := &subscriberv1.GetFriendProfileRequest{
		FriendId: friendID,
	}

	resp, err := s.GRPCClient.GetFriendProfile(ctx, req)
	if err != nil {
		return entity.Friend{}, s.mapSubscriberError(err)
	}

	friendData := resp.GetFriend()
	if friendData == nil {
		s.log.Error("invalid response from subscriber service",
			slog.String("op", op),
			slog.String("issue", "friend profile data is nil"),
		)
		return entity.Friend{}, domain.ErrInvalidResponse
	}

	return entity.Friend{
		ID:           friendData.GetUserId(),
		Name:         friendData.GetName(),
		FriendsSince: friendData.GetFriendsSince().AsTime(),
	}, nil
}

func (s *SubscriberClient) RemoveFriend(ctx context.Context, friendID string) (err error) {
	req := &subscriberv1.DeleteFriendRequest{
		FriendId: friendID,
	}

	_, err = s.GRPCClient.DeleteFriend(ctx, req)
	if err != nil {
		return s.mapSubscriberError(err)
	}

	return nil
}

func (s *SubscriberClient) InviteFriend(ctx context.Context, invitation entity.FriendInvite) (inviteID string, err error) {
	const op = "grpc.client.subscriber.InviteFriend"

	req := &subscriberv1.SendFriendInviteRequest{
		ToUserId: invitation.ToUserID,
		Message:  &invitation.Message,
	}

	resp, err := s.GRPCClient.SendFriendInvite(ctx, req)
	if err != nil {
		return "", s.mapSubscriberError(err)
	}

	inviteID = resp.GetInviteId()
	if inviteID == "" {
		s.log.Error("invalid response from subscriber service",
			slog.String("op", op),
			slog.String("issue", "invite id is empty"),
		)
		return "", domain.ErrInvalidResponse
	}

	return inviteID, nil
}

func (s *SubscriberClient) GetFriendInvites(ctx context.Context, direction, status, pageToken *string) (invites []entity.FriendInvite, nextPageToken string, err error) {
	req := &subscriberv1.GetAllFriendInvitesRequest{}

	if direction != nil {
		if dir, ok := subscriberv1.RequestDirection_value[*direction]; ok {
			dirEnum := subscriberv1.RequestDirection(dir)
			req.Direction = &dirEnum
		}
	}

	if status != nil {
		if st, ok := subscriberv1.RequestStatus_value[*status]; ok {
			statusEnum := subscriberv1.RequestStatus(st)
			req.Status = &statusEnum
		}
	}

	if pageToken != nil {
		req.PageToken = *pageToken
	}

	resp, err := s.GRPCClient.GetAllFriendInvites(ctx, req)
	if err != nil {
		return nil, "", s.mapSubscriberError(err)
	}

	invitesData := resp.GetInvites()
	if invitesData == nil {
		return nil, "", domain.ErrNoFriendInvites
	}

	invites = make([]entity.FriendInvite, len(invitesData))
	for i, inviteData := range invitesData {
		invites[i] = entity.FriendInvite{
			ID:         inviteData.GetId(),
			FromUserID: inviteData.GetFromUserId(),
			ToUserID:   inviteData.GetToUserId(),
			Status:     string(inviteData.GetStatus()),
			Message:    inviteData.GetMessage(),
			CreatedAt:  inviteData.GetCreatedAt().AsTime(),
			UpdatedAt:  inviteData.GetUpdatedAt().AsTime(),
		}
	}

	return invites, nextPageToken, nil
}

func (s *SubscriberClient) AcceptFriendInvite(ctx context.Context, inviteID string) (err error) {
	req := &subscriberv1.AcceptFriendInviteRequest{
		InviteId: inviteID,
	}

	_, err = s.GRPCClient.AcceptFriendInvite(ctx, req)
	if err != nil {
		return s.mapSubscriberError(err)
	}

	return nil
}

func (s *SubscriberClient) DeclineFriendInvite(ctx context.Context, inviteID string) (err error) {
	req := &subscriberv1.DeclineFriendInviteRequest{
		InviteId: inviteID,
	}

	_, err = s.GRPCClient.DeclineFriendInvite(ctx, req)
	if err != nil {
		return s.mapSubscriberError(err)
	}

	return nil
}

func (s *SubscriberClient) mapSubscriberError(err error) error {
	// TODO: implement this
	return nil
}
