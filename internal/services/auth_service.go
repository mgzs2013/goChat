package services

import (
	"database/sql"
	"errors"
	"fmt"
	"goChat/internal/database"
	"goChat/internal/models"
	"goChat/pkg"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(username, password string) (int64, string, error) {
	userID, role, err := getUserCredentials(username, password) // Assume this is implemented
	if err != nil {
		return 0, "", errors.New("invalid credentials")
	}
	return userID, role, nil
}

func GenerateToken(secretKey string, ID int64, username, role string) (string, string, error) {
	// Fetch the secret from the environment
	secret := os.Getenv("JWT_SECRET")
	if len(secret) == 0 {
		return "", "", fmt.Errorf("JWT secret is missing")
	}

	log.Printf("[DEBUG] JWT_SECRET being used for validation: %s", secret)

	// Step 1: Dynamically fetch user ID from the database
	query := `SELECT id FROM users WHERE username = $1`
	err := database.Pool.QueryRow(query, username).Scan(&ID)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch user ID: %w", err)
	}

	// Debug log to verify correct ID
	log.Printf("[DEBUG] Fetched user ID for username '%s': %d", username, ID)
	// Generate Access Token
	accessTokenClaims := jwt.MapClaims{
		"customKey": "mySuperSecretKey",
		"id":        ID,                                   // User ID
		"username":  username,                             // Username
		"role":      role,                                 // User role
		"exp":       time.Now().Add(2 * time.Hour).Unix(), // Token expires in 2 hours
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	signedToken, err := accessToken.SignedString([]byte(secret))

	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}
	log.Printf("[DEBUG] Access Token Claims: %v", accessTokenClaims)
	log.Printf("[DEBUG] Attempting to insert token for user ID: %d", ID)

	// Generate Refresh Token
	refreshTokenClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(), // Refresh token expires in 7 days
		IssuedAt:  time.Now().Unix(),
		Issuer:    "auth_service",
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Insert Access Token into the Database
	query2 := `INSERT INTO access_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err = database.Pool.Exec(query2, ID, signedToken, time.Now().Add(24*time.Hour))
	if err != nil {
		return "", "", fmt.Errorf("failed to store access token in database: %v", err)
	}

	log.Printf("Generated and stored token with exp: %d, current time: %d", time.Now().Add(24*time.Hour).Unix(), time.Now().Unix())
	log.Printf("[DEBUG] Access Token String: %v", signedToken)

	return signedToken, refreshTokenString, nil
}

func ValidateToken(signedToken string) (jwt.MapClaims, error) {

	secret := []byte(os.Getenv("JWT_SECRET"))
	log.Printf("[DEBUG] Loaded JWT_SECRET: %s", string(secret))
	if len(secret) == 0 {
		log.Println("JWT_SECRET is not set or empty")
		return nil, fmt.Errorf("JWT secret is missing")
	}

	// Parse the token
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		log.Printf("[DEBUG] Token header: %+v", token.Header)                // Log header details
		log.Printf("[DEBUG] Token claims before parsing: %+v", token.Claims) // Log claims if available

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		log.Printf("[ERROR] Error parsing token: %v", err)
		log.Printf("[DEBUG] Error during token parsing: %v", err)
		return nil, err
	}

	// Extract claims and validate token
	log.Println("[DEBUG] Attempting to extract claims from token")
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	log.Printf("[DEBUG] Extracted Claims: %v", claims)
	if !ok {
		log.Printf("[ERROR] Claims type assertion failed. Token Claims: %v", parsedToken.Claims)
		log.Println("Failed to extract claims from token")
		return nil, fmt.Errorf("invalid token claims")
	}

	if !parsedToken.Valid {
		log.Println("Token is invalid")
		return nil, fmt.Errorf("invalid token")
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		log.Println("[ERROR] Expiration (exp) claim missing or invalid")
		return nil, fmt.Errorf("expiration claim invalid")
	}
	if int64(exp) < time.Now().Unix() {
		log.Println("[ERROR] Token has expired")
		return nil, fmt.Errorf("token expired")
	}
	log.Printf("[DEBUG] Token expiration time: %d, Current time: %d", int64(exp), time.Now().Unix())
	log.Println("Token validated successfully")
	return claims, nil

}

func GetUserByID(userID int64) (*models.User, error) {
	// Fetch full user details from the database
	var user models.User
	err := database.Pool.QueryRow(
		"SELECT id, username, role FROM users WHERE id = $1",
		userID).Scan(
		&user.ID,
		&user.Username,
		&user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("Error fetching user by ID: %v", err)
		return nil, err
	}
	log.Println("user retreived by GetUserByID function:", user)
	return &user, nil
}

func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("Bcrypt password comparison failed:", err) // Debugging line
		return false
	}
	return true
}

func getUserCredentials(username, password string) (int64, string, error) {
	var userID int64
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
func GenerateRefreshToken(userID int64) (string, error) {
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
func ValidateRefreshToken(refreshToken string) (int64, error) {
	var userID int64
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
