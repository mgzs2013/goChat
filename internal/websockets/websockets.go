package websockets

import (
	"log"
	"net/http"

	"goChat/config"

	"goChat/internal/middleware"

	"github.com/gorilla/websocket"
)

type User struct {
	Username string
	// Add other relevant fields as needed
}

// Upgrader upgrades HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// For now, allow all connections. In production, restrict this.
		return true
	},
}

// Message defines the structure for messages exchanged between clients.
type Message struct {
	Username  string `json:"username"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Typing    bool   `json:"typing"`
}

// A map to track active clients.
var clients = make(map[*websocket.Conn]*User)

// A broadcast channel for incoming messages.
var broadcast = make(chan Message)

func HandleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.ValidateJWT(r, config.JwtSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	// Register new client
	hub.RegisterClient(ws, &User{Username: claims.Username})

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("ReadJSON error: %v", err)
			hub.RemoveClient(ws)
			break
		}
		msg.Username = claims.Username
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
