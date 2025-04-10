package models

import "time"

type Message struct {
	ID        int64     `json:"id"`
	SenderID  int64     `json:"sender_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
