package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"goChat/internal/database"
	"goChat/internal/handlers"
	"goChat/internal/middleware"
	"goChat/internal/websockets"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	} else {
		fmt.Println("Successfully loaded .env file")
	}

	// Check if the secret is loaded
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is empty! Check your .env file")
	} else {
		fmt.Println("Loaded JWT_SECRET:", secret)
	}

	// Initialize database
	dbConnectionString := os.Getenv("DATABASE_URL")
	if dbConnectionString == "" {
		log.Fatal("DATABASE_URL is empty! Check your .env file")
	}

	database.InitDB(dbConnectionString)

	// Create a new Hub for managing WebSocket clients and messages
	hub := websockets.NewHub()

	// Start handling WebSocket messages in a separate goroutine
	go websockets.HandleMessages(hub)

	// Setup router and define routes
	r := http.NewServeMux()

	// Authentication route
	r.HandleFunc("/login", handlers.HandleLogin)

	// Messages route
	r.Handle("/messages/history", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetMessageHistoryHandler)))

	// WebSocket route
	r.Handle("/ws", middleware.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		websockets.HandleConnections(hub, w, r)
	})))

	// Role-based routes
	r.Handle("/admin", middleware.RoleMiddleware("admin", http.HandlerFunc(handlers.AdminHandler)))
	r.Handle("/editor", middleware.RoleMiddleware("editor", http.HandlerFunc(handlers.EditorHandler)))

	// Start the server
	log.Println("Server started on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
