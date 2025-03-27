package pkg

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// JwtSecret defined in config.go

func GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
