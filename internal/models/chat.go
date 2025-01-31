package models

import "time"

type ChatType string

const (
	ChatTypePrivate ChatType = "private"
	ChatTypeGroup   ChatType = "group"
)

type Chat struct {
	ID        string    `json:"id"`
	Type      ChatType  `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatMember struct {
	ChatID   string    `json:"chat_id"`
	UserID   string    `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
}
