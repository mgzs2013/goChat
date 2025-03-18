package handlers

import (
	"encoding/json"

	"goChat/internal/repository"
	"goChat/internal/services" // Import your service layer
	"log"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		log.Printf("Invalid input: %+v", req)
		http.Error(w, "Invalid input: Username and Password are required", http.StatusBadRequest)
		return
	}

	// Use service layer for authentication
	ID, role, err := services.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		log.Println("Login failed:", err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := services.GenerateToken(ID, req.Username, role)
	if err != nil {
		log.Println("Error generating tokens:", err)
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	log.Printf("Login successful for username: %s", req.Username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshTokenHandler generates new access and refresh tokens
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if refreshToken == "" {
		http.Error(w, "Refresh token missing", http.StatusUnauthorized)
		return
	}

	// Validate the refresh token
	userID, err := services.ValidateRefreshToken(refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Fetch the user role from the repository
	role, err := repository.GetUserRole(userID)
	if err != nil {
		http.Error(w, "Failed to fetch user role", http.StatusInternalServerError)
		return
	}

	// Fetch username for the user (optional based on your requirements)
	username, err := repository.GetUserRole(userID)
	if err != nil {
		http.Error(w, "Failed to fetch username", http.StatusInternalServerError)
		return
	}

	// Generate new tokens (access token and refresh token)
	accessToken, newRefreshToken, err := services.GenerateToken(userID, username, role)
	if err != nil {
		http.Error(w, "Failed to generate new tokens", http.StatusInternalServerError)
		return
	}

	// Send the new tokens in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"accessToken":  accessToken,
		"refreshToken": newRefreshToken,
	})
}
