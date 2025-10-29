package handler

import "errors"

var (
	ErrInternalServerError                    = errors.New("internal server error")
	ErrBadRequest                             = errors.New("bad request")
	ErrInvalidRequest                         = errors.New("invalid request")
	ErrValidationFailed                       = errors.New("request validation failed")
	ErrFailedToExtractRefreshTokenFromCookies = errors.New("failed to extract refresh token from cookies")
	ErrVerificationTokenRequired              = errors.New("verification token is required")
	ErrPasswordResetTokenRequired             = errors.New("password reset token is required")
	ErrChatIDRequired                         = errors.New("chat ID is required")
	ErrFriendIDRequired                       = errors.New("friend ID is required")
	ErrInviteIDRequired                       = errors.New("invite ID is required")
)
