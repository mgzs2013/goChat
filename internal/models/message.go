package models

import "time"

type Message struct {
	ID          int64     `json:"id"`
	SenderID    int64     `json:"sender_id"`
	RecipientID int64     `json:"recipient_id,omitempty"` // Optional for group chats
	ChatRoom    string    `json:"chat_room,omitempty"`    // Optional for private messages
	Content     string    `json:"content"`
	Timestamp   time.Time `json:"timestamp"`
}
