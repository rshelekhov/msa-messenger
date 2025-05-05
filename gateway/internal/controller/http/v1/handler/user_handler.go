package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain"
)

func (h *Handler) GetOwnProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.GetOwnProfile"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		user, err := h.userUsecase.GetOwnProfile(ctx)
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				handleError(w, r, log, err, http.StatusNotFound)
				return
			}
			handleInternalError(w, r, log, err)
			return
		}

		resp := toGetOwnProfileResponse(user)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}

func (h *Handler) UpdateOwnProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.UpdateOwnProfile"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		req := &UpdateOwnProfileRequest{}
		if err := h.decodeAndValidateRequest(w, r, log, req); err != nil {
			return
		}

		user := toUpdateUser(req)

		var err error
		user.ID, err = h.userUsecase.UpdateOwnProfile(ctx, user)
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				handleError(w, r, log, err, http.StatusNotFound)
				return
			}
			handleInternalError(w, r, log, err)
			return
		}

		resp := toUpdateUserResponse(user)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}

func (h *Handler) DeleteOwnProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.DeleteOwnProfile"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		err := h.userUsecase.DeleteOwnProfile(ctx)
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				handleError(w, r, log, err, http.StatusNotFound)
				return
			}
			handleInternalError(w, r, log, err)
			return
		}

		render.Status(r, http.StatusOK)
	}
}

func (h *Handler) SearchUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.SearchUsers"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		params := &GetUsersSearchParams{}
		if err := render.Decode(r, params); err != nil {
			err = fmt.Errorf("failed to decode request: %w", err)
			handleBadRequestError(w, r, log, err)
			return
		}

		users, err := h.userUsecase.SearchUsers(ctx, params.Query)
		if err != nil {
			handleInternalError(w, r, log, err)
			return
		}

		resp := toSearchUsersResponse(users)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}
