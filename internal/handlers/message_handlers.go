package handlers

import (
	"encoding/json"
	"fmt"

	"goChat/internal/models"
	"goChat/internal/services"
	"net/http"
)

func CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	var msg models.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the service function to store the message
	newID, err := services.StoreMessage(msg)
	if err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Message saved with ID: %d", newID)))
}

func GetMessageHistoryHandler(w http.ResponseWriter, r *http.Request) {
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
