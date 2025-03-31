package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

func main() {
	// Replace with your WebSocket server URL
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21LZXkiOiJteVN1cGVyU2VjcmV0S2V5IiwiZXhwIjoxNzQzMzg0NDI0LCJpZCI6Miwicm9sZSI6ImFkbWluIiwidXNlcm5hbWUiOiJhZG1pbnVzZXIifQ.BFg7m5Ce9jtL6OWpADJIgP5ZyY4vMZbnJlV5UdkjwuU" // Replace with the token from the CURL login response
	url := fmt.Sprintf("ws://localhost:8080/ws?accessToken=%s", accessToken)

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
		os.Exit(1)
	}
	defer conn.Close()

	log.Println("WebSocket connection established")

	// Send a login confirmation message
	msg := Message{
		Type: "LOGIN_CONFIRMATION",
		ID:   "2", // Replace with the user ID from the login response
	}
	err = conn.WriteJSON(msg)
	if err != nil {
		log.Fatal("Error sending message:", err)
	}

	// Read messages from the server
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received message from server: %s", message)
	}
}
