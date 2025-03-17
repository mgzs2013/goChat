package websockets

import (
	"goChat/internal/models"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients    map[*websocket.Conn]*models.User
	Broadcast  chan models.Message
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*websocket.Conn]*models.User),
		Broadcast:  make(chan models.Message),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

func (hub *Hub) RegisterClient(conn *websocket.Conn, user *models.User) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.Clients[conn] = user
}

func (hub *Hub) RemoveClient(conn *websocket.Conn) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	delete(hub.Clients, conn)
}

func (hub *Hub) BroadcastMessage(msg models.Message) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	for conn := range hub.Clients {
		err := conn.WriteJSON(msg)
		if err != nil {
			conn.Close()
			delete(hub.Clients, conn)
		}
	}
}
