package handler

import (
	"errors"
	"log/slog"
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
			log.Warn("failed to decode request", slog.String("error", err.Error()))
			handleBadRequestError(w, r, ErrInvalidRequest)
			return
		}

		friends, nextPageToken, err := h.subscriberUsecase.GetFriends(ctx, params.PageToken)
		if err != nil {
			log.Error("failed to get friends", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
			log.Warn("friend ID is required")
			handleBadRequestError(w, r, ErrFriendIDRequired)
			return
		}

		friend, err := h.subscriberUsecase.GetFriend(ctx, friendID)
		if err != nil {
			if errors.Is(err, domain.ErrFriendNotFound) {
				log.Warn("failed to get friend", slog.String("error", err.Error()))
				handleError(w, r, domain.ErrFriendNotFound, http.StatusNotFound)
				return
			}

			log.Error("failed to get friend", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
			log.Warn("friend ID is required")
			handleBadRequestError(w, r, ErrFriendIDRequired)
			return
		}

		err := h.subscriberUsecase.RemoveFriend(ctx, friendID)
		if err != nil {
			if errors.Is(err, domain.ErrFriendNotFound) {
				log.Warn("failed to remove friend", slog.String("error", err.Error()))
				handleError(w, r, domain.ErrFriendNotFound, http.StatusNotFound)
				return
			}

			log.Error("failed to remove friend", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
			log.Warn("friend ID is required")
			handleBadRequestError(w, r, ErrFriendIDRequired)
			return
		}

		req := &InviteFriendRequest{}
		if err := h.decodeAndValidateRequest(w, r, log, req); err != nil {
			return
		}

		invitation := toInviteFriend(friendID, req)

		invitationID, err := h.subscriberUsecase.InviteFriend(ctx, invitation)
		if err != nil {
			log.Error("failed to invite friend", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
			log.Warn("failed to decode request", slog.String("error", err.Error()))
			handleBadRequestError(w, r, ErrInvalidRequest)
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
			log.Error("failed to get friend invites", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
			log.Warn("invite ID is required")
			handleBadRequestError(w, r, ErrInviteIDRequired)
			return
		}

		if err := h.subscriberUsecase.AcceptFriendInvite(ctx, inviteID); err != nil {
			if errors.Is(err, domain.ErrFriendInviteNotFound) {
				log.Warn("failed to accept friend invite", slog.String("error", err.Error()))
				handleError(w, r, domain.ErrFriendInviteNotFound, http.StatusNotFound)
				return
			}

			log.Error("failed to accept friend invite", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
			log.Warn("invite ID is required")
			handleBadRequestError(w, r, ErrInviteIDRequired)
			return
		}

		if err := h.subscriberUsecase.DeclineFriendInvite(ctx, inviteID); err != nil {
			if errors.Is(err, domain.ErrFriendInviteNotFound) {
				log.Warn("failed to decline friend invite", slog.String("error", err.Error()))
				handleError(w, r, domain.ErrFriendInviteNotFound, http.StatusNotFound)
				return
			}

			log.Error("failed to decline friend invite", slog.String("error", err.Error()))
			handleInternalError(w, r)
			return
		}

		render.Status(r, http.StatusOK)
	}
}
