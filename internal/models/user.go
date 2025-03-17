package models

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
	Role         string `json:"role,omitempty"`
}
