package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rshelekhov/jwtauth"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/entity"
)

type tokenSender struct{}

func NewTokenSender() *tokenSender {
	return &tokenSender{}
}

func (s *tokenSender) Send(w http.ResponseWriter, r *http.Request, response any, tokens entity.UserTokens, httpStatus int) error {
	if isMobileRequest(r) {
		return s.sendTokensToMobile(w, response, httpStatus)
	}

	return s.sendTokensToWeb(w, response, tokens, httpStatus)
}

func (s *tokenSender) sendTokensToWeb(w http.ResponseWriter, response any, tokens entity.UserTokens, httpStatus int) error {
	const op = "token_manager.sendTokensToWeb"

	s.setRefreshTokenCookie(w, tokens)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf("%s: failed to encode response: %w", op, err)
	}

	return nil
}

func (s *tokenSender) setRefreshTokenCookie(w http.ResponseWriter, tokens entity.UserTokens) {
	http.SetCookie(w, &http.Cookie{
		Name:     jwtauth.RefreshTokenKey,
		Value:    tokens.RefreshToken,
		Domain:   tokens.Domain,
		Path:     tokens.Path,
		Expires:  tokens.ExpiresAt,
		HttpOnly: tokens.HttpOnly,
	})
}

func (s *tokenSender) sendTokensToMobile(w http.ResponseWriter, response any, httpStatus int) error {
	const op = "token_manager.sendTokensToMobile"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf("%s: failed to encode response: %w", op, err)
	}

	return nil
}
