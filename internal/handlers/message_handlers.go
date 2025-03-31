package handlers

import (
	"encoding/json"
	"log"

	"goChat/internal/models"
	"goChat/internal/services"
	"goChat/internal/websockets"
	"net/http"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	// Ensure only POST method is accepted
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var message models.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate message (add more validation as needed)
	if message.Content == "" || message.SenderID == 0 {
		http.Error(w, "Invalid message data", http.StatusBadRequest)
		return
	}

	// Store the message in the database
	messageID, err := services.StoreMessage(message)
	if err != nil {
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}

	// Broadcast the message (if using WebSocket)
	websockets.NewHub().BroadcastMessage(message)

	// Respond with the created message ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": messageID})
}

// HandleIncomingMessage processes incoming WebSocket messages
func HandleIncomingMessage(hub *websockets.Hub, messageData []byte) {
	var message models.Message

	// Unmarshal the JSON message
	if err := json.Unmarshal(messageData, &message); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	// Store the message in the database
	messageID, err := services.StoreMessage(message)
	if err != nil {
		log.Printf("Error storing message: %v", err)
		return
	}

	// Optional: Add the generated message ID to the message
	message.ID = messageID

	// Broadcast the message to all connected clients
	hub.Broadcast <- message
}

// GetChatHistory retrieves chat history for a specific chat room or between users
func GetMessageHistory(w http.ResponseWriter, r *http.Request) {
	senderID := r.URL.Query().Get("senderId")
	recipientID := r.URL.Query().Get("recipientId")
	chatRoom := r.URL.Query().Get("chatRoom")

	if (senderID == "" || recipientID == "") && chatRoom == "" {
		http.Error(w, "Invalid parameters: senderId and recipientId or chatRoom must be provided", http.StatusBadRequest)
		return
	}

	messages, err := services.GetMessageHistory(senderID, recipientID, chatRoom)
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

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
