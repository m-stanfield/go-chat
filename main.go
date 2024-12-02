package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Upgrader to upgrade HTTP connections to WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for simplicity
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected, starting to send messages")

	// Send a JSON message every second for 20 seconds
	for i := 1; i <= 20; i++ {
		message := map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"message":   fmt.Sprintf("Message %d", i),
			"counter":   i,
		}

		jsonMessage, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal message: %v\n", err)
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			log.Printf("Failed to send message: %v\n", err)
			break
		}

		log.Printf("Sent: %s\n", jsonMessage)

		// Wait for 1 second before sending the next message
		time.Sleep(1 * time.Second)
	}

	log.Println("Finished sending messages, closing connection")
}

func main() {
	http.HandleFunc("/", handleWebSocket)

	port := "8081"
	log.Printf("WebSocket server started on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

