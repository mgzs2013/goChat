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
	// Validate the JWT and get claims
	claims, err := middleware.ValidateJWT(r, config.JwtSecret)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch the user ID from the database
	var userID int64
	err = database.Pool.QueryRow("SELECT id FROM users WHERE username = $1", claims.Username).Scan(&userID)
	if err != nil {
		log.Printf("Failed to fetch user ID for username %s: %v", claims.Username, err)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Register the client with the hub
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	hub.RegisterClient(ws, &models.User{ID: userID, Username: claims.Username})

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("ReadJSON error: %v", err)
			hub.RemoveClient(ws)
			break
		}
		// Set the SenderID to the fetched user ID
		msg.SenderID = userID
		hub.Broadcast <- msg
	}
}

// HandleMessages listens on the broadcast channel and sends messages to all clients.
func HandleMessages(hub *Hub) {
	for {
		msg := <-hub.Broadcast
		hub.BroadcastMessage(msg)
	}
}
