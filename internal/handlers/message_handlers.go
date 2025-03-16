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
