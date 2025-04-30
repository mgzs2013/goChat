package handlers

import (
	"encoding/json"
	"goChat/internal/services"
	"strings"

	"log"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("RegisterUserHanlder invoked!")

	var req RegisterRequest
	if r.Method != http.MethodPost {
		RegisterRespondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "This is not a valid method!"})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[ERROR] Cannot decode request body: %v", err)
		RegisterRespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request body",
			"details": "Ensure the request body is valid JSON with 'username' and 'password' fields.",
		})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Password == "" {
		RegisterRespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid input",
			"details": "Username and password cannot be empty.",
		})
		return
	}

	log.Println("username and password from the client:", req.Username, req.Password)

	err := services.RegisterUser(req.Username, req.Password)
	if err != nil {
		log.Println("Error Registering User:", err)
		RegisterRespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "User was not registered properly"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})

}

func RegisterRespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
