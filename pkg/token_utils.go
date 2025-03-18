package pkg

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// JwtSecret defined in config.go
func GenerateToken(ID int64, username, role string) (string, string, error) {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) == 0 {
		return "", "", fmt.Errorf("JWT secret is missing")
	}

	// Generate Access Token
	accessTokenClaims := jwt.MapClaims{
		"id":       ID,                                    // User ID as an integer
		"username": username,                              // Username as a string
		"role":     role,                                  // Role as a string (e.g., "admin", "editor")
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Access token expires in 2 hours
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	// Generate Refresh Token (typically with minimal claims)
	refreshTokenClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(), // Refresh token expires in 7 days
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	log.Printf("Generated token with exp: %d, current time: %d", time.Now().Add(2*time.Hour).Unix(), time.Now().Unix())

	return accessTokenString, refreshTokenString, nil
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		log.Println("JWT_SECRET is not set or empty")
		return nil, fmt.Errorf("JWT secret is missing")
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check for unexpected signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, err
	}

	// Extract claims and validate token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Failed to extract claims from token")
		return nil, fmt.Errorf("invalid token claims")
	}

	if !token.Valid {
		log.Println("Token is invalid")
		return nil, fmt.Errorf("invalid token")
	}

	log.Println("Token validated successfully")
	return claims, nil
}

func GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
