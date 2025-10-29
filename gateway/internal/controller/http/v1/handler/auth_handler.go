package handler

import (
	"errors"
	"log/slog"
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
			if errors.Is(err, domain.ErrUserAlreadyExists) || errors.Is(err, domain.ErrInvalidRequest) {
				log.Warn("registration failed", slog.String("error", err.Error()))
				handleError(w, r, err, http.StatusConflict)
				return
			}

			log.Error("failed to register user", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
			log.Error("failed to send tokens to client", slog.String("error", err.Error()))
			handleInternalError(w, r)
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

		log = log.With(
			slog.String("email", userCredentials.Email),
			slog.String("user_agent", userDevice.UserAgent),
		)

		userID, tokens, err := h.authUsecase.Login(ctx, userCredentials, userDevice)
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrInvalidCredentials) {
				log.Warn("authentication failed", slog.String("error", err.Error()))
				handleError(w, r, err, http.StatusUnauthorized)
				return
			}

			log.Error("failed to login user", slog.String("error", err.Error()))
			handleInternalError(w, r)
			return
		}

		var resp any
		if isMobileRequest(r) {
			resp = &LoginResponseMobile{
				Id:           &userID,
				AccessToken:  &tokens.AccessToken,
				RefreshToken: &tokens.RefreshToken,
			}
		} else {
			resp = &LoginResponseWeb{
				Id:          &userID,
				AccessToken: &tokens.AccessToken,
			}
		}

		if err := h.tokenSender.Send(w, r, resp, tokens, http.StatusOK); err != nil {
			log.Error("failed to send tokens to client", slog.String("error", err.Error()))
			handleInternalError(w, r)
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

		log = log.With(
			slog.String("user_agent", userDevice.UserAgent),
		)

		if err := h.authUsecase.Logout(ctx, userDevice); err != nil {
			if errors.Is(err, domain.ErrInvalidRequest) {
				log.Warn("failed to logout user", slog.String("error", err.Error()))
				handleError(w, r, domain.ErrInvalidRequest, http.StatusBadRequest)
				return
			} else if errors.Is(err, domain.ErrUserDeviceNotRegisteredOrAlreadyLoggedOut) {
				log.Warn("failed to logout user", slog.String("error", err.Error()))
				handleError(w, r, domain.ErrUserDeviceNotRegisteredOrAlreadyLoggedOut, http.StatusUnauthorized)
				return
			}

			log.Error("failed to logout user", slog.String("error", err.Error()))
			handleInternalError(w, r)
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

		render.Status(r, http.StatusOK)
		render.JSON(w, r, &LogoutResponse{
			Message: "Logged out successfully",
		})
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
				log.Error("failed to extract refresh token from cookies", slog.String("error", err.Error()))
				handleBadRequestError(w, r, ErrFailedToExtractRefreshTokenFromCookies)
				return
			}
		}

		userDevice := toUserDevice(r)

		tokens, err := h.authUsecase.RefreshToken(ctx, refreshToken, userDevice)
		if err != nil {
			errorMappings := map[error]int{
				domain.ErrInvalidCredentials:                        http.StatusBadRequest,
				domain.ErrSessionNotFound:                           http.StatusUnauthorized,
				domain.ErrSessionExpired:                            http.StatusUnauthorized,
				domain.ErrUserDeviceNotRegisteredOrAlreadyLoggedOut: http.StatusUnauthorized,
			}
			if handleMappedError(w, r, err, log, "failed to refresh tokens", errorMappings) {
				return
			}

			log.Error("failed to refresh tokens", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
			log.Error("failed to send tokens to client", slog.String("error", err.Error()))
			handleInternalError(w, r)
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
			log.Error("failed to get JWKS", slog.String("error", err.Error()))
			handleInternalError(w, r)
			return
		}

		jwksResp := toJWKSResponse(jwks)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, jwksResp)
	}
}

func (h *Handler) VerifyEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.VerifyEmail"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		verificationToken := r.URL.Query().Get("token")
		if verificationToken == "" {
			log.Warn("verification token is required")
			handleBadRequestError(w, r, ErrVerificationTokenRequired)
			return
		}

		err := h.authUsecase.VerifyEmail(ctx, verificationToken)
		if err != nil {
			if errors.Is(err, domain.ErrTokenExpiredEmailResent) || errors.Is(err, domain.ErrInvalidRequest) {
				log.Warn("failed to verify email", slog.String("error", err.Error()))
				handleError(w, r, err, http.StatusBadRequest)
				return
			}

			log.Error("failed to verify email", slog.String("error", err.Error()))
			handleInternalError(w, r)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, &VerifyEmailResponse{
			Message: "Email verified successfully",
		})
	}
}

func (h *Handler) ResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.ResetPassword"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		req := &ResetPasswordRequest{}
		if err := h.decodeAndValidateRequest(w, r, log, req); err != nil {
			return
		}

		err := h.authUsecase.ResetPassword(ctx, req.Email)
		if err != nil {
			errorMappings := map[error]int{
				domain.ErrTokenExpiredEmailResent: http.StatusBadRequest,
				domain.ErrInvalidRequest:          http.StatusBadRequest,
				domain.ErrUserNotFound:            http.StatusNotFound,
			}
			if handleMappedError(w, r, err, log, "failed to reset password", errorMappings) {
				return
			}

			log.Error("failed to reset password", slog.String("error", err.Error()))
			handleInternalError(w, r)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, &ResetPasswordResponse{
			Message: "Password reset email sent",
		})
	}
}

func (h *Handler) ChangePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.ChangePassword"

		ctx := r.Context()
		log := h.logWithReqID(ctx, op)

		passwordResetToken := r.URL.Query().Get("token")
		if passwordResetToken == "" {
			log.Warn("password reset token is required")
			handleBadRequestError(w, r, ErrPasswordResetTokenRequired)
			return
		}

		req := &ChangePasswordRequest{}
		if err := h.decodeAndValidateRequest(w, r, log, req); err != nil {
			return
		}

		err := h.authUsecase.ChangePassword(ctx, passwordResetToken, req.UpdatedPassword)
		if err != nil {
			errorMappings := map[error]int{
				domain.ErrTokenExpiredEmailResent:   http.StatusBadRequest,
				domain.ErrInvalidRequest:            http.StatusBadRequest,
				domain.ErrNoPasswordChangesDetected: http.StatusBadRequest,
			}
			if handleMappedError(w, r, err, log, "failed to change password", errorMappings) {
				return
			}

			log.Error("failed to change password", slog.String("error", err.Error()))
			handleInternalError(w, r)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, &ChangePasswordResponse{
			Message: "Password changed successfully",
		})
	}
}
