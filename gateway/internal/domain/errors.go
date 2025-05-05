package domain

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrFriendNotFound       = errors.New("friend not found")
	ErrFriendInviteNotFound = errors.New("friend invite not found")
	ErrChatNotFound         = errors.New("chat not found")
)
