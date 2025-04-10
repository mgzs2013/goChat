package handlers

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"goChat/internal/models"
	"goChat/internal/services"
	"goChat/internal/websockets"
	"net/http"
)

type MessageRequest struct {
	ID        int64     `json:"id"`
	SenderID  int64     `json:"sender_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type MessageResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Assuming you have a Message struct for the database
type Message struct {
	ID        int64  `db:"id"`
	SenderID  int64  `db:"sender_id"`
	Content   string `db:"content"`
	Timestamp time.Time
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] CreateMessage Invoked")

	var req MessageRequest
	// Ensure only POST method is accepted
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[ERROR] Cannot decode request body: %v", err)
		MessageRespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request body",
			"details": "Ensure the request body is valid JSON with 'message' fields.",
		})
		return
	}

	req.Content = strings.TrimSpace(req.Content)

	// Validate message (add more validation as needed)
	if req.Content == "" || req.SenderID == 0 {
		http.Error(w, "Invalid message data", http.StatusBadRequest)
		return
	}

	log.Println("The message being sent through the socket and the sender_id:", req.Content, req.SenderID)
	// Store the message in the database
	err := services.StoreMessage(req.SenderID, req.Content, req.Timestamp)
	if err != nil {
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}

	// Broadcast the message to all connected clients
	websockets.NewHub().Broadcast <- models.Message{

		SenderID:  req.SenderID,
		Content:   req.Content,
		Timestamp: req.Timestamp,
	}
	// Respond with the created message ID

	w.WriteHeader(http.StatusCreated)

}

// HandleIncomingMessage processes incoming WebSocket messages
func HandleIncomingMessage(hub *websockets.Hub, messageData []byte) {
	log.Println("[DEBUG] HandleIncomingMessage invoked")
	// Parse the request body
	var message struct {
		ID        int64     `json:"id"`
		SenderID  int64     `json:"sender_id"`
		Content   string    `json:"content"`
		Timestamp time.Time `json:"timestamp"`
	}
	// Unmarshal the JSON message
	if err := json.Unmarshal(messageData, &message); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	// Store the message in the database
	err := services.StoreMessage(message.SenderID, message.Content, message.Timestamp)
	if err != nil {
		log.Printf("Error storing message: %v", err)
		return
	}

	// Optional: Add the generated message ID to the message

	// Broadcast the message to all connected clients
	websockets.NewHub().Broadcast <- models.Message{

		SenderID:  message.SenderID,
		Content:   message.Content,
		Timestamp: message.Timestamp,
	}
}

func MessageRespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// GetChatHistory retrieves chat history for a specific chat room or between users
// func GetMessageHistory(w http.ResponseWriter, r *http.Request) {
// 	senderID := r.URL.Query().Get("senderId")
// 	recipientID := r.URL.Query().Get("recipientId")
// 	chatRoom := r.URL.Query().Get("chatRoom")

// 	if (senderID == "" || recipientID == "") && chatRoom == "" {
// 		http.Error(w, "Invalid parameters: senderId and recipientId or chatRoom must be provided", http.StatusBadRequest)
// 		return
// 	}

// 	messages, err := services.GetMessageHistory(senderID, recipientID, chatRoom)
// 	if err != nil {
// 		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(messages)
// }

// HandleRecentMessages retrieves recent messages based on criteria
func HandleRecentMessages(w http.ResponseWriter, r *http.Request) {
	// You can implement additional logic for fetching recent messages
	// This could include pagination, limit, etc.

	// Example basic implementation
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "50" // Default limit
	}

	// TODO: Implement logic to fetch recent messages with limit
	// This might require adding a new method in message_service.go
	messages, err := services.GetRecentMessages(limit)
	if err != nil {
		http.Error(w, "Failed to retrieve recent messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
