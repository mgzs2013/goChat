package websockets

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients   map[*websocket.Conn]*User
	Broadcast chan Message
}

func NewHub() *Hub {
	return &Hub{
		Clients:   make(map[*websocket.Conn]*User),
		Broadcast: make(chan Message),
	}
}

func (h *Hub) RegisterClient(ws *websocket.Conn, user *User) {
	h.Clients[ws] = user
}

func (h *Hub) RemoveClient(ws *websocket.Conn) {
	delete(h.Clients, ws)
}

func (h *Hub) BroadcastMessage(msg Message) {
	for client := range h.Clients {
		err := client.WriteJSON(msg)
		if err != nil {
			client.Close()
			h.RemoveClient(client)
		}
	}
}
