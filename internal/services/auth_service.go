package services

import (
	"errors"
	"fmt"
	"goChat/internal/database"

	"goChat/pkg"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(username, password string) (int, string, error) {
	userID, role, err := getUserCredentials(username, password) // Assume this is implemented
	if err != nil {
		return 0, "", errors.New("invalid credentials")
	}
	return userID, role, nil
}

func GenerateToken(userID int, role string) (string, string, error) {
	accessToken, err := pkg.GenerateToken(fmt.Sprintf("user_%d", userID), role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("Bcrypt password comparison failed:", err) // Debugging line
		return false
	}
	return true
}

func getUserCredentials(username, password string) (int, string, error) {
	var userID int
	var role string
	var storedHash string

	// Fetch user credentials from DB
	err := database.Pool.QueryRow("SELECT id, password_hash, role FROM users WHERE username = $1", username).Scan(&userID, &storedHash, &role)
	if err != nil {
		log.Println("Error fetching user credentials:", err) // Debugging line
		return 0, "", err
	}

	log.Printf("Retrieved User - ID: %d, Role: %s, Hash: %s\n", userID, role, storedHash) // Debugging line

	// Verify password
	if !CheckPassword(password, storedHash) {
		log.Println("Password mismatch for user:", username) // Debugging line
		return 0, "", errors.New("invalid credentials")
	}

	log.Println("User authenticated successfully:", username) // Debugging line
	return userID, role, nil
}

// GenerateRefreshToken creates and stores a new refresh token in the database
func GenerateRefreshToken(userID int) (string, error) {
	// Use GenerateRandomToken to generate the refresh token
	refreshToken, err := pkg.GenerateRandomToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	expiresAt := time.Now().Add(24 * time.Hour) // Set expiration to 24 hours

	// Remove existing refresh tokens for the user
	_, err = database.Pool.Exec(
		"DELETE FROM refresh_tokens WHERE user_id = $1", userID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to remove old refresh tokens: %w", err)
	}

	// Insert the new refresh token into the database
	_, err = database.Pool.Exec(
		"INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, refreshToken, expiresAt,
	)
	if err != nil {
		return "", fmt.Errorf("failed to store new refresh token: %w", err)
	}

	return refreshToken, nil
}

// ValidateRefreshToken checks the validity of a refresh token in the database
func ValidateRefreshToken(refreshToken string) (int, error) {
	var userID int
	var expiresAt time.Time

	// Query the refresh token from the database
	err := database.Pool.QueryRow(
		"SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1", refreshToken,
	).Scan(&userID, &expiresAt)

	if err != nil {
		return 0, fmt.Errorf("refresh token not found or invalid")
	}

	// Check if the refresh token has expired
	if time.Now().After(expiresAt) {
		return 0, fmt.Errorf("refresh token has expired")
	}

	return userID, nil
}
