package services

import (
	"goChat/internal/database"
	"goChat/internal/models"
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

func GetMessageHistory(senderID, recipientID, chatRoom string) ([]models.Message, error) {
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

	rows, err := database.Pool.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.RecipientID, &msg.ChatRoom, &msg.Content, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
