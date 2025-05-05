package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain"
)

func (h *Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.Register"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		req := &RegisterRequest{}
		if err := h.decodeAndValidateRequest(w, r, log, req); err != nil {
			return
		}

		userCredentials := toRegisterUserCredentials(req)
		userDevice := toUserDevice(r)

		userID, tokens, err := h.authUsecase.Register(ctx, userCredentials, userDevice)
		if err != nil {
			if errors.Is(err, domain.ErrUserAlreadyExists) {
				handleError(w, r, log, err, http.StatusConflict)
				return
			}
			handleInternalError(w, r, log, err)
			return
		}

		var resp any
		if isMobileRequest(r) {
			resp = &RegisterResponseMobile{
				AccessToken:  &tokens.AccessToken,
				RefreshToken: &tokens.RefreshToken,
				Id:           &userID,
				Name:         &userCredentials.Name,
				Email:        &userCredentials.Email,
			}
		} else {
			resp = &RegisterResponseWeb{
				AccessToken: &tokens.AccessToken,
				Id:          &userID,
				Name:        &userCredentials.Name,
				Email:       &userCredentials.Email,
			}
		}

		if err := h.tokenSender.Send(w, r, resp, tokens, http.StatusCreated); err != nil {
			handleInternalError(w, r, log, err)
			return
		}
	}
}

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.Login"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		req := &LoginRequest{}
		if err := h.decodeAndValidateRequest(w, r, log, req); err != nil {
			return
		}

		userCredentials := toLoginUserCredentials(req)
		userDevice := toUserDevice(r)

		tokens, err := h.authUsecase.Login(ctx, userCredentials, userDevice)
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				err = fmt.Errorf("invalid credentials")
				handleError(w, r, log, err, http.StatusUnauthorized)
				return
			}
			handleInternalError(w, r, log, err)
			return
		}

		var resp any
		if isMobileRequest(r) {
			resp = &LoginResponseMobile{
				AccessToken:  &tokens.AccessToken,
				RefreshToken: &tokens.RefreshToken,
			}
		} else {
			resp = &LoginResponseWeb{
				AccessToken: &tokens.AccessToken,
			}
		}

		if err := h.tokenSender.Send(w, r, resp, tokens, http.StatusOK); err != nil {
			handleInternalError(w, r, log, err)
			return
		}
	}
}

func (h *Handler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.Logout"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		userDevice := toUserDevice(r)

		if err := h.authUsecase.Logout(ctx, userDevice); err != nil {
			handleInternalError(w, r, log, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Domain:   "",
			Path:     "/",
			Expires:  time.Now().Add(-time.Hour * 24),
			HttpOnly: true,
			MaxAge:   -1,
		})

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("logged out"))
		if err != nil {
			handleInternalError(w, r, log, err)
			return
		}
	}
}

func (h *Handler) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.RefreshToken"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		var refreshToken string
		var err error
		if isMobileRequest(r) {
			req := &RefreshTokenRequestMobile{}
			if err := h.decodeAndValidateRequest(w, r, log, req); err != nil {
				return
			}
			refreshToken = req.RefreshToken
		} else {
			refreshToken, err = h.tokenManager.ExtractRefreshTokenFromCookies(r)
			if err != nil {
				handleBadRequestError(w, r, log, err)
				return
			}
		}

		userDevice := toUserDevice(r)

		tokens, err := h.authUsecase.RefreshToken(ctx, refreshToken, userDevice)
		if err != nil {
			handleInternalError(w, r, log, err)
			return
		}

		var resp any
		if isMobileRequest(r) {
			resp = &RefreshTokenResponseMobile{
				AccessToken:  &tokens.AccessToken,
				RefreshToken: &tokens.RefreshToken,
			}
		} else {
			resp = &RefreshTokenResponseWeb{
				AccessToken: &tokens.AccessToken,
			}
		}

		if err := h.tokenSender.Send(w, r, resp, tokens, http.StatusOK); err != nil {
			handleInternalError(w, r, log, err)
			return
		}
	}
}

func (h *Handler) GetJWKS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.GetJWKS"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		jwks, err := h.authUsecase.GetJWKS(ctx)
		if err != nil {
			handleInternalError(w, r, log, err)
			return
		}

		jwksResp := toJWKSResponse(jwks)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, jwksResp)
	}
}
