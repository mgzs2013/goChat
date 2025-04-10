package websockets

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"goChat/internal/models"
	"goChat/internal/services"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
)

// Upgrader upgrades HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 10 * time.Second,
}

// A map to track active clients.
var clients = make(map[*websocket.Conn]*models.User)

// A broadcast channel for incoming messages.
var broadcast = make(chan models.Message)

func HandleWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] HandleWebsocket invoked")

	// Extract the token from the query parameters
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		log.Println("[ERROR] Access token is missing in query parameters")
		http.Error(w, "Access token required", http.StatusBadRequest)
		return
	}
	log.Printf("[DEBUG] Extracted Token: %s", tokenString)

	// Validate the token
	claims, err := services.ValidateToken(tokenString)
	if err != nil {
		log.Println("[ERROR] Token validation failed:", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	user, err := getUserFromClaims(claims)
	if err != nil {
		log.Printf("[ERROR] Could not retrieve user from claims: %v", err)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Configure upgrader with ping/pong handlers and connection close handler
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// Implement your origin checking logic here
		return true
	}

	// Upgrade the connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	// Register the connection with the hub
	hub.Register <- conn

	// Track the user with the connection
	hub.mu.Lock()
	hub.Clients[conn] = user
	hub.mu.Unlock()

	// Set up ping/pong handlers to keep connection alive
	conn.SetPingHandler(func(message string) error {
		err := conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
		return err
	})

	// Create done channel to manage connection lifecycle
	done := make(chan struct{})

	// Start a goroutine to handle incoming messages
	go func() {
		defer func() {
			hub.Unregister <- conn
			close(done)
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("[ERROR] WebSocket read error: %v", err)
				}
				break
			}

			msg := models.Message{

				SenderID:  user.ID,
				Content:   string(message),
				Timestamp: time.Now(),
			}

			log.Println("[DEBUG] message through websocket:", msg)
			err = services.StoreMessage(msg.SenderID, msg.Content, msg.Timestamp)
			if err != nil {
				log.Printf("[ERROR] Failed to store message: %v", err)
				continue // Handle the error as needed
			}
			// Broadcast the message
			hub.Broadcast <- msg
		}
	}()

	// Wait for connection to be closed
	<-done
}

// Helper function to get user from claims (implement based on your models)
func getUserFromClaims(claims jwt.MapClaims) (*models.User, error) {
	log.Printf("[DEBUG] JWT_SECRET being used for validation: %+v", claims)

	// Check if claims are nil
	if claims == nil {
		return nil, fmt.Errorf("token claims are nil")
	}

	// Extract user ID from claims
	// Use type conversion and check to handle potential float64 from JSON parsing
	var userID int64
	switch v := claims["id"].(type) {
	case float64:
		userID = int64(v)
	case int64:
		userID = v
	case string:
		// If ID is stored as string, convert it
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID format")
		}
		userID = id
	default:
		return nil, fmt.Errorf("user ID not found or in invalid format")
	}

	// Fetch user from database or user service
	user, err := services.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}

	return user, nil
}

// HandleMessages listens on the broadcast channel and sends messages to all clients.
func HandleMessages(hub *Hub) {
	for {
		msg := <-hub.Broadcast
		hub.BroadcastMessage(msg)
	}
}
