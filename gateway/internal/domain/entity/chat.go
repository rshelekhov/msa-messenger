package entity

import "time"

type Chat struct {
	ID        string
	UserID    string
	FriendID  string
	CreatedAt time.Time
	Messages  []Message
}

type Message struct {
	ID        string
	SenderID  string
	Content   string
	CreatedAt time.Time
}
