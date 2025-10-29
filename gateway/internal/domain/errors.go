package domain

import "errors"

var (
	// Common errors
	ErrInvalidResponse = errors.New("invalid response")
	ErrInvalidRequest  = errors.New("invalid request")

	// SSO (common) errors
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyTaken = errors.New("email already taken")

	// SSO (auth)errors
	ErrInvalidCredentials                        = errors.New("invalid credentials")
	ErrSessionExpired                            = errors.New("session expired")
	ErrSessionNotFound                           = errors.New("session not found")
	ErrUserAlreadyExists                         = errors.New("user already exists")
	ErrTokenExpiredEmailResent                   = errors.New("token expired, new email sent")
	ErrFailedToSendVerificationEmail             = errors.New("failed to send verification email")
	ErrFailedToSendResetPasswordEmail            = errors.New("failed to send reset password email")
	ErrUserDeviceNotRegisteredOrAlreadyLoggedOut = errors.New("user device not registered or already logged out")
	ErrNoUsersFound                              = errors.New("no users found")

	// SSO (user) errors
	ErrCurrentPasswordRequired    = errors.New("current password required")
	ErrNoEmailChangesDetected     = errors.New("no email changes detected")
	ErrNoPasswordChangesDetected  = errors.New("no password changes detected")
	ErrNoNameChangesDetected      = errors.New("no name changes detected")
	ErrPasswordsDoNotMatch        = errors.New("passwords do not match")
	ErrCurrentPasswordIsIncorrect = errors.New("current password is incorrect")

	// Chat errors
	ErrChatNotFound = errors.New("chat not found")

	// Subscriber errors
	ErrFriendNotFound       = errors.New("friend not found")
	ErrFriendInviteNotFound = errors.New("friend invite not found")
	ErrNoFriends            = errors.New("no friends")
	ErrNoFriendInvites      = errors.New("no friend invites")
)
