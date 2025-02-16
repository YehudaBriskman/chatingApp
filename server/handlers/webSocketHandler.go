package handlers

import (
	"chatingApp/services"
	"log"
	"net/http"
	"strconv"
	"sync"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections for chat rooms.
type WebSocketHandler struct {
	RoomService *services.RoomService
	Clients     map[int]map[*websocket.Conn]bool // roomID -> set of WebSocket connections
	Mutex       sync.Mutex
}

// NewWebSocketHandler creates a new WebSocketHandler instance.
func NewWebSocketHandler(service *services.RoomService) *WebSocketHandler {
	return &WebSocketHandler{
		RoomService: service,
		Clients:     make(map[int]map[*websocket.Conn]bool),
	}
}

// Upgrader for WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// HandleWebSocketConnection manages WebSocket connections per room.
func (h *WebSocketHandler) HandleWebSocketConnection(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("roomID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	// userID := c.GetInt("user_id") // Fetch user ID from middleware
	// if !h.RoomService.IsUserInRoom(roomID, userID) {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "User not in room"})
	// 	return
	// }

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("❌ WebSocket Upgrade Failed:", err)
		return
	}
	defer conn.Close()

	h.Mutex.Lock()
	if _, exists := h.Clients[roomID]; !exists {
		h.Clients[roomID] = make(map[*websocket.Conn]bool)
	}
	h.Clients[roomID][conn] = true
	h.Mutex.Unlock()

	// log.Printf("✅ WebSocket Connection Established (Room: %d, User: %d)\n", roomID, userID)

	for {
		var msg struct {
			UserID  int    `json:"user_id"`
			Content string `json:"content"`
		}

		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("❌ WebSocket Read Error:", err)
			break
		}

		// log.Printf("📨 Message Received (Room: %d, User: %d): %s\n", roomID, msg.UserID, msg.Content)

		// Save message in database
		// err := h.RoomService.AddMessageToRoom(roomID, msg.UserID, msg.Content)
		// if err != nil {
		// 	log.Println("❌ Failed to save message:", err)
		// 	continue
		// }

		// Broadcast message to all clients in the room
		h.broadcastMessage(roomID, msg.UserID, msg.Content)
	}

	// Cleanup on disconnect
	h.Mutex.Lock()
	delete(h.Clients[roomID], conn)
	h.Mutex.Unlock()
	// log.Printf("❌ WebSocket Disconnected (Room: %d, User: %d)\n", roomID, userID)
}

// Broadcast message to all WebSocket clients in a room.
func (h *WebSocketHandler) broadcastMessage(roomID, senderID int, message string) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	for client := range h.Clients[roomID] {
		err := client.WriteJSON(gin.H{"room_id": roomID, "user_id": senderID, "content": message})
		if err != nil {
			log.Println("❌ WebSocket Write Error:", err)
			client.Close()
			delete(h.Clients[roomID], client)
		}
	}
}
