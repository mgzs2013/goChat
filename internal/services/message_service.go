package services

import (
	"fmt"
	"goChat/internal/database"

	"goChat/internal/models"
	"log"
	"strconv"
	"time"
)

func StoreMessage(SenderID int64, Content string, Timestamp time.Time) error {
	log.Println("[DEBUG] StoreMessage is being invoked")

	// Use double quotes around "Message" to match the case
	query := `INSERT INTO "Message" (sender_id, content, timestamp) VALUES ($1, $2, $3)`
	_, err := database.Pool.Exec(query, SenderID, Content, Timestamp)
	if err != nil {
		return fmt.Errorf("failed to store message in database: %v", err)
	}
	return nil
}

func GetRecentMessages(limit string) ([]models.Message, error) {
	query := `
		SELECT id, sender_id, content, timestamp
		FROM messages
		ORDER BY timestamp DESC
		LIMIT $1
	`

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 50 // Default limit
	}

	rows, err := database.Pool.Query(query, limitInt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.Content, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// func GetMessageHistory(senderID, recipientID, chatRoom string) ([]models.Message, error) {
// 	var query string
// 	var args []interface{}

// 	if chatRoom != "" {
// 		query = `
//             SELECT id, sender_id, content, timestamp
//             FROM messages
//             WHERE chat_room = $1
//             ORDER BY timestamp
//         `
// 		args = append(args, chatRoom)
// 	} else {
// 		query = `
//             SELECT id, sender_id, content, timestamp
//             FROM messages
//             WHERE (sender_id = $1 AND recipient_id = $2)
//                OR (sender_id = $2 AND recipient_id = $1)
//             ORDER BY timestamp
//         `
// 		args = append(args, senderID, recipientID)
// 	}

// 	rows, err := database.Pool.Query(query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var messages []models.Message
// 	for rows.Next() {
// 		var msg models.Message
// 		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.RecipientID, &msg.ChatRoom, &msg.Content, &msg.Timestamp); err != nil {
// 			return nil, err
// 		}
// 		messages = append(messages, msg)
// 	}

// 	return messages, nil
// }
