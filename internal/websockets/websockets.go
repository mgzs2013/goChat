package websockets

import (
	"log"
	"net/http"

	"goChat/pkg"

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

func HandleWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] HandleWebsocket invoked")

	// Extract the token from the query parameters
	tokenString := r.URL.Query().Get("accessToken")
	if tokenString == "" {
		log.Println("[ERROR] Access token is missing in query parameters")
		http.Error(w, "Access token required", http.StatusBadRequest)
		return
	}
	log.Printf("[DEBUG] Extracted Token: %s", tokenString)

	// Validate the token
	claims, err := pkg.ValidateToken(tokenString)
	if err != nil {
		log.Println("[ERROR] Token validation failed:", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	log.Println("[DEBUG] Token validated successfully. Claims:", claims)

	// Upgrade the connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] WebSocket upgrade failed:", err)
		return
	}
	log.Println("[DEBUG] WebSocket connection established")
	defer conn.Close()

	// WebSocket handling logic here
}

// HandleMessages listens on the broadcast channel and sends messages to all clients.
func HandleMessages(hub *Hub) {
	for {
		msg := <-hub.Broadcast
		hub.BroadcastMessage(msg)
	}
}
