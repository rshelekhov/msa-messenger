package handler

import (
	"net/http"
	"strings"

	"github.com/rshelekhov/jwtauth"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/entity"
)

func toRegisterUserCredentials(request *PostAuthRegisterJSONRequestBody) entity.UserCredentials {
	return entity.UserCredentials{
		Email:    request.Email,
		Name:     request.Name,
		Password: request.Password,
	}
}

func toLoginUserCredentials(request *PostAuthLoginJSONRequestBody) entity.UserCredentials {
	return entity.UserCredentials{
		Email:    request.Email,
		Password: request.Password,
	}
}

func toUserDevice(r *http.Request) entity.UserDevice {
	return entity.UserDevice{
		UserAgent: r.UserAgent(),
		IP:        strings.Split(r.RemoteAddr, ":")[0],
	}
}

func toJWKSResponse(jwks entity.JWKS) *jwtauth.JWKSResponse {
	keys := make([]jwtauth.JWK, len(jwks.Keys))
	for i, key := range jwks.Keys {
		keys[i] = jwtauth.JWK{
			Kty: key.Kty,
			Use: key.Use,
			Alg: key.Alg,
			Kid: key.Kid,
			N:   key.N,
			E:   key.E,
		}
	}
	return &jwtauth.JWKSResponse{
		Keys: keys,
		TTL:  jwks.TTL,
	}
}

func toUpdateUser(request *PatchUsersMeJSONRequestBody) entity.User {
	return entity.User{
		Name:  *request.Name,
		Email: *request.Email,
	}
}

func toGetOwnProfileResponse(user entity.User) *UserResponse {
	return &UserResponse{
		Id:    &user.ID,
		Name:  &user.Name,
		Email: &user.Email,
	}
}

func toUpdateUserResponse(user entity.User) *UserResponse {
	return &UserResponse{
		Id:    &user.ID,
		Name:  &user.Name,
		Email: &user.Email,
	}
}

func toSearchUsersResponse(users []entity.User) *[]UserResponse {
	resp := make([]UserResponse, len(users))
	for i, user := range users {
		resp[i] = UserResponse{
			Id:    &user.ID,
			Name:  &user.Name,
			Email: &user.Email,
		}
	}
	return &resp
}

func toGetFriendsResponse(friends []entity.Friend, nextPageToken string) *FriendsResponse {
	friendsResp := make([]FriendResponse, len(friends))
	for i, friend := range friends {
		friendsResp[i] = FriendResponse{
			Id:           &friend.ID,
			Name:         &friend.Name,
			FriendsSince: &friend.FriendsSince,
		}
	}
	return &FriendsResponse{
		Friends:       &friendsResp,
		NextPageToken: &nextPageToken,
	}
}

func toGetFriendResponse(friend entity.Friend) *FriendResponse {
	return &FriendResponse{
		Id:           &friend.ID,
		Name:         &friend.Name,
		FriendsSince: &friend.FriendsSince,
	}
}

func toInviteFriend(friendID string, request *PostFriendsIdInviteJSONRequestBody) entity.FriendInvite {
	var friendInvite entity.FriendInvite

	friendInvite.ToUserID = friendID

	if request.Message != nil {
		friendInvite.Message = *request.Message
	}

	return friendInvite
}

func toInviteFriendResponse(invitationID string) *InviteFriendResponse {
	return &InviteFriendResponse{
		Id: &invitationID,
	}
}

func toFriendInvitesResponse(invites []entity.FriendInvite, nextPageToken string) *InvitesResponse {
	resp := make([]FriendInvite, len(invites))
	for i, invite := range invites {
		resp[i] = FriendInvite{
			Id:         &invite.ID,
			FromUserId: &invite.FromUserID,
			ToUserId:   &invite.ToUserID,
			Status:     (*FriendInviteStatus)(&invite.Status),
			Message:    &invite.Message,
			CreatedAt:  &invite.CreatedAt,
			UpdatedAt:  &invite.UpdatedAt,
		}
	}
	return &InvitesResponse{
		Invites:       &resp,
		NextPageToken: &nextPageToken,
	}
}

func toGetChatsResponse(chats []entity.Chat, nextPageToken string) *ChatsResponse {
	resp := make([]Chat, len(chats))
	for i, chat := range chats {
		resp[i] = Chat{
			Id:        &chat.ID,
			UserId:    &chat.UserID,
			FriendId:  &chat.FriendID,
			CreatedAt: &chat.CreatedAt,
		}
	}
	return &ChatsResponse{
		Chats:         &resp,
		NextPageToken: &nextPageToken,
	}
}

func toGetChatResponse(chat entity.Chat, nextPageToken string) *MessagesResponse {
	chatResp := Chat{
		Id:        &chat.ID,
		UserId:    &chat.UserID,
		FriendId:  &chat.FriendID,
		CreatedAt: &chat.CreatedAt,
	}

	messagesResp := make([]MessageResponse, len(chat.Messages))
	for i, message := range chat.Messages {
		messagesResp[i] = MessageResponse{
			Id:        &message.ID,
			Content:   &message.Content,
			CreatedAt: &message.CreatedAt,
		}
	}

	return &MessagesResponse{
		Chat:          &chatResp,
		Messages:      &messagesResp,
		NextPageToken: &nextPageToken,
	}
}

func toSendMessageResponse(message entity.Message) *MessageResponse {
	return &MessageResponse{
		Content:   &message.Content,
		CreatedAt: &message.CreatedAt,
		Id:        &message.ID,
		SenderId:  &message.SenderID,
	}
}
