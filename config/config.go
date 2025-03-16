package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JwtSecret []byte

func init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set JwtSecret from the environment
	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(JwtSecret) == 0 {
		log.Fatal("JWT_SECRET is not set in the environment")
	}
}
