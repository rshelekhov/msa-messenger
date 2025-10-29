package handler

import (
	"errors"
	"log/slog"
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
				log.Warn("failed to get own profile", slog.String("error", err.Error()))
				handleError(w, r, domain.ErrUserNotFound, http.StatusNotFound)
				return
			}

			log.Error("failed to get own profile", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
			errorMappings := map[error]int{
				domain.ErrUserNotFound:               http.StatusNotFound,
				domain.ErrInvalidRequest:             http.StatusBadRequest,
				domain.ErrCurrentPasswordRequired:    http.StatusBadRequest,
				domain.ErrNoEmailChangesDetected:     http.StatusBadRequest,
				domain.ErrNoPasswordChangesDetected:  http.StatusBadRequest,
				domain.ErrNoNameChangesDetected:      http.StatusBadRequest,
				domain.ErrPasswordsDoNotMatch:        http.StatusBadRequest,
				domain.ErrCurrentPasswordIsIncorrect: http.StatusBadRequest,
				domain.ErrEmailAlreadyTaken:          http.StatusBadRequest,
			}
			if handleMappedError(w, r, err, log, "failed to update own profile", errorMappings) {
				return
			}

			log.Error("failed to update own profile", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
				log.Warn("failed to delete own profile", slog.String("error", err.Error()))
				handleError(w, r, domain.ErrUserNotFound, http.StatusNotFound)
				return
			}

			log.Error("failed to delete own profile", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
			log.Warn("failed to decode request", slog.String("error", err.Error()))
			handleBadRequestError(w, r, ErrInvalidRequest)
			return
		}

		users, err := h.userUsecase.SearchUsers(ctx, params.Query)
		if err != nil {
			log.Error("failed to search users", slog.String("error", err.Error()))
			handleInternalError(w, r)
			return
		}

		resp := toSearchUsersResponse(users)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}
