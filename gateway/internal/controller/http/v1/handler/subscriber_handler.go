package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain"
)

func (h *Handler) GetFriends() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.subscriber.GetFriends"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		params := &GetFriendsParams{}
		if err := render.Decode(r, params); err != nil {
			err = fmt.Errorf("failed to decode request: %w", err)
			handleBadRequestError(w, r, log, err)
			return
		}

		friends, nextPageToken, err := h.subscriberUsecase.GetFriends(ctx, params.PageToken)
		if err != nil {
			handleInternalError(w, r, log, err)
			return
		}

		resp := toGetFriendsResponse(friends, nextPageToken)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}

func (h *Handler) GetFriend() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.subscriber.GetFriend"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		friendID := chi.URLParam(r, "id")
		if friendID == "" {
			handleBadRequestError(w, r, log, errors.New("friend ID is required"))
			return
		}

		friend, err := h.subscriberUsecase.GetFriend(ctx, friendID)
		if err != nil {
			if errors.Is(err, domain.ErrFriendNotFound) {
				handleError(w, r, log, err, http.StatusNotFound)
				return
			}
			handleInternalError(w, r, log, err)
			return
		}

		resp := toGetFriendResponse(friend)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}

func (h *Handler) RemoveFriend() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.subscriber.RemoveFriend"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		friendID := chi.URLParam(r, "id")
		if friendID == "" {
			handleBadRequestError(w, r, log, errors.New("friend ID is required"))
			return
		}

		err := h.subscriberUsecase.RemoveFriend(ctx, friendID)
		if err != nil {
			if errors.Is(err, domain.ErrFriendNotFound) {
				handleError(w, r, log, err, http.StatusNotFound)
				return
			}
			handleInternalError(w, r, log, err)
			return
		}

		render.Status(r, http.StatusOK)
	}
}

func (h *Handler) InviteFriend() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.subscriber.InviteFriend"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		friendID := chi.URLParam(r, "id")
		if friendID == "" {
			handleBadRequestError(w, r, log, errors.New("friend ID is required"))
			return
		}

		req := &InviteFriendRequest{}
		if err := h.decodeAndValidateRequest(w, r, log, req); err != nil {
			return
		}

		invitation := toInviteFriend(friendID, req)

		invitationID, err := h.subscriberUsecase.InviteFriend(ctx, invitation)
		if err != nil {
			handleInternalError(w, r, log, err)
			return
		}

		resp := toInviteFriendResponse(invitationID)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}

func (h *Handler) GetFriendInvites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.subscriber.GetFriendInvites"
		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		params := &GetFriendInvitesParams{}
		if err := render.Decode(r, params); err != nil {
			err = fmt.Errorf("failed to decode request: %w", err)
			handleBadRequestError(w, r, log, err)
			return
		}

		var direction, status *string
		if params.Direction != nil {
			s := string(*params.Direction)
			direction = &s
		}
		if params.Status != nil {
			s := string(*params.Status)
			status = &s
		}

		invites, nextPageToken, err := h.subscriberUsecase.GetFriendInvites(
			ctx,
			direction,
			status,
			params.PageToken,
		)
		if err != nil {
			handleInternalError(w, r, log, err)
			return
		}

		resp := toFriendInvitesResponse(invites, nextPageToken)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}

func (h *Handler) AcceptFriendInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.subscriber.AcceptFriendInvite"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		inviteID := chi.URLParam(r, "id")
		if inviteID == "" {
			handleBadRequestError(w, r, log, errors.New("invite ID is required"))
			return
		}

		if err := h.subscriberUsecase.AcceptFriendInvite(ctx, inviteID); err != nil {
			if errors.Is(err, domain.ErrFriendInviteNotFound) {
				handleError(w, r, log, err, http.StatusNotFound)
				return
			}
			handleInternalError(w, r, log, err)
			return
		}

		render.Status(r, http.StatusOK)
	}
}

func (h *Handler) DeclineFriendInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.subscriber.DeclineFriendInvite"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		inviteID := chi.URLParam(r, "id")
		if inviteID == "" {
			handleBadRequestError(w, r, log, errors.New("invite ID is required"))
			return
		}

		if err := h.subscriberUsecase.DeclineFriendInvite(ctx, inviteID); err != nil {
			if errors.Is(err, domain.ErrFriendInviteNotFound) {
				handleError(w, r, log, err, http.StatusNotFound)
				return
			}
			handleInternalError(w, r, log, err)
			return
		}

		render.Status(r, http.StatusOK)
	}
}
