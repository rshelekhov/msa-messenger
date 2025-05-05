package entity

import "time"

type (
	Friend struct {
		ID           string
		Name         string
		FriendsSince time.Time
	}

	FriendInvite struct {
		ID         string
		FromUserID string
		ToUserID   string
		Status     string
		Message    string
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}
)
