package services

import (
	"encoding/json"
	"fmt"
	"goChat/internal/database"
	"goChat/internal/models"
	"net/http"
)

func StoreMessage(msg models.Message) (int, error) {
	query := `
        INSERT INTO messages (sender_id, recipient_id, chat_room, content, timestamp) 
        VALUES ($1, $2, $3, $4, DEFAULT) RETURNING id
    `
	var newID int
	err := database.Pool.QueryRow(query, msg.SenderID, msg.RecipientID, msg.ChatRoom, msg.Content).Scan(&newID)
	return newID, err
}

func GetMessageHistory(w http.ResponseWriter, r *http.Request) {
	senderID := r.URL.Query().Get("senderId")
	recipientID := r.URL.Query().Get("recipientId")
	chatRoom := r.URL.Query().Get("chatRoom")

	// Validate required parameters
	if (senderID == "" || recipientID == "") && chatRoom == "" {
		http.Error(w, "Invalid parameters: senderId and recipientId or chatRoom must be provided", http.StatusBadRequest)
		return
	}

	// Build query
	var query string
	var args []interface{}
	if chatRoom != "" {
		query = `
			SELECT id, sender_id, recipient_id, chat_room, content, timestamp
			FROM messages
			WHERE chat_room = $1
			ORDER BY timestamp
		`
		args = append(args, chatRoom)
	} else {
		query = `
			SELECT id, sender_id, recipient_id, chat_room, content, timestamp
			FROM messages
			WHERE (sender_id = $1 AND recipient_id = $2)
			   OR (sender_id = $2 AND recipient_id = $1)
			ORDER BY timestamp
		`
		args = append(args, senderID, recipientID)
	}

	// Fetch data
	rows, err := database.Pool.Query(query, args...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database query failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.RecipientID, &msg.ChatRoom, &msg.Content, &msg.Timestamp); err != nil {
			http.Error(w, "Failed to parse messages", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	// Respond with messages
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
