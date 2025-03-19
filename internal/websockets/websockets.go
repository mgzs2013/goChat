package websockets

import (
	"log"
	"net/http"

	"goChat/config"

	"goChat/internal/database"
	"goChat/internal/middleware"
	"goChat/internal/models"

	"github.com/gorilla/websocket"
)

// Upgrader upgrades HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// For now, allow all connections. In production, restrict this.
		return true
	},
}

// A map to track active clients.
var clients = make(map[*websocket.Conn]*models.User)

// A broadcast channel for incoming messages.
var broadcast = make(chan models.Message)

func HandleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.ValidateJWT(r, config.JwtSecret)
	if err != nil {
		log.Printf("JWT validation failed: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var userID int64
	err = database.Pool.QueryRow("SELECT id FROM users WHERE username = $1", claims.Username).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
		return
	}
	defer hub.RemoveClient(ws)

	user := &models.User{ID: userID, Username: claims.Username}
	hub.RegisterClient(ws, user)
}

// HandleMessages listens on the broadcast channel and sends messages to all clients.
func HandleMessages(hub *Hub) {
	for {
		msg := <-hub.Broadcast
		hub.BroadcastMessage(msg)
	}
}
