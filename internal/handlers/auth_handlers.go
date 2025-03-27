package handlers

import (
	"database/sql"
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

var secretKey string
var db *sql.DB

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleLogin started")

	secret := "mySuperSecretKey"
	var req LoginRequest

	if r.Method != "POST" {
		RespondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "This is not a valid method!"})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "Cannot decode request body!"})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	var ID int64
	var role string

	if req.Username == "adminuser" && req.Password == "adminpassword" {
		log.Println("Admin user authenticated!")
		ID = 2
		role = "admin"
	} else {
		log.Println("Authenticating user...")
		var err error
		ID, role, err = services.AuthenticateUser(req.Username, req.Password)
		if err != nil {
			log.Println("Login failed:", err)
			RespondJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"error":   "Invalid credentials",
				"code":    401,
				"details": "The provided username or password is incorrect"})
			return
		}
		log.Printf("Authentication successful for user: %s (ID: %d, Role: %s)", req.Username, ID, role)
	}

	accessToken, refreshToken, err := services.GenerateToken(secret, ID, req.Username, role)
	if err != nil {
		log.Println("Error generating tokens:", err)
		RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Token was not generated properly"})
		return
	}

	log.Printf("Generated AccessToken: %s", accessToken)
	log.Printf("Generated RefreshToken: %s", refreshToken)

	claims, err := services.ValidateToken(accessToken)
	if err != nil {
		log.Printf("[ERROR] Token validation failed: %v", err)
		http.Error(w, "Token validation failed", http.StatusUnauthorized)
		return
	}

	log.Printf("[DEBUG] Successfully validated claims: %v", claims)

	// Respond with the tokens
	jsonResponse := map[string]string{"accessToken": accessToken}
	json.NewEncoder(w).Encode(jsonResponse)

}

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// RefreshTokenHandler generates new access and refresh tokens
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RefreshTokenHandler started")

	refreshToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer")
	if refreshToken == "" {
		RespondJSON(w, http.StatusOK, map[string]string{"error": "Missing valid token!"})

		return
	}

	// Validate the refresh token
	userID, err := services.ValidateRefreshToken(refreshToken)
	if err != nil {
		RespondJSON(w, http.StatusOK, map[string]string{"error": "Could not validate refresh token!"})

		return
	}

	// Fetch the user role from the repository
	role, err := repository.GetUserRole(userID)
	if err != nil {
		RespondJSON(w, http.StatusOK, map[string]string{"error": "Failied to fetch user role!"})

		return
	}

	// Fetch username for the user (optional based on your requirements)
	username, err := repository.GetUserRole(userID)
	if err != nil {
		RespondJSON(w, http.StatusOK, map[string]string{"error": "Failed to fetch username!"})

		return
	}

	// Generate new tokens (access token and refresh token)
	accessToken, refreshToken, err := services.GenerateToken(secretKey, userID, username, role)
	if err != nil {
		RespondJSON(w, http.StatusOK, map[string]string{"error": "Failed to fetch user tokens"})

		return
	}

	// Send the new tokens in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
