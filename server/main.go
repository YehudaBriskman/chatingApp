package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

type Message struct {
	Type     string `json:"type"`     
	Username string `json:"username"` 
	Text     string `json:"text"`     
	UserID   string `json:"userId"`   
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]string)
var broadcast = make(chan Message)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	userID := uuid.New().String()
	clients[conn] = userID

	fmt.Println("New user connected:", r.RemoteAddr, "User ID:", userID)

	broadcast <- Message{
		Type:     "join",
		Username: "User-" + userID[:5],
		Text:     "has joined the chat!",
		UserID:   userID,
	}

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("User disconnected:", r.RemoteAddr, "User ID:", userID)
			delete(clients, conn)
			broadcast <- Message{
				Type:     "leave",
				Username: msg.Username,
				Text:     "has left the chat.",
				UserID:   userID,
			}
			break
		}

		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Error sending message:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	fmt.Println("Chat server started on ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
