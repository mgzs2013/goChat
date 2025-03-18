package websockets

import (
	"log"
	"net/http"

	"goChat/config"

	"goChat/internal/database"
	"goChat/internal/middleware"
	"goChat/internal/models"
	"goChat/internal/websockets"

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
		log.Printf("JWT validation failed: %v", err)
		http.Error(w, "Unauthorized: invalid or expired token", http.StatusUnauthorized)
		return
	}
	log.Printf("JWT validated successfully for username: %s", claims.Username)

	// Fetch the user ID from the database
	var userID int64
	err = database.Pool.QueryRow("SELECT id FROM users WHERE username = $1", claims.Username).Scan(&userID)
	if err != nil {
		log.Printf("Failed to fetch user ID for username %s: %v", claims.Username, err)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	log.Printf("User ID fetched for username %s: %d", claims.Username, userID)

	// Upgrade the HTTP connection to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		http.Error(w, "Failed to upgrade WebSocket connection", http.StatusInternalServerError)
		return
	}
	defer func() {
		hub.RemoveClient(ws)
		log.Printf("Client disconnected: %s", claims.Username)
	}()

	log.Println("WebSocket connection established!")

	// Register the client with the hub
	user := &models.User{ID: userID, Username: claims.Username}
	if err := websockets.Hub.RegisterClient(ws, user); err != nil {
		log.Printf("Failed to register client: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send acknowledgment to the client
	ackMessage := map[string]string{"status": "connected", "message": "Welcome!"}
	if err := ws.WriteJSON(ackMessage); err != nil {
		log.Printf("Failed to send acknowledgment: %v", err)
	}

	// Handle incoming messages
	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket closed unexpectedly: %v", err)
			} else {
				log.Printf("ReadJSON error: %v", err)
			}
			break
		}
		log.Printf("Message received: %+v", msg)

		// Set the SenderID and broadcast the message
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
