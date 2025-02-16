package models

import "time"

// Room represents a chat room or group where users can send messages.
type Room struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"` // Optional room description
	CreatedBy   int       `json:"created_by"`            // User ID of the creator
	RoomAdmins  []int     `json:"room_admins"`           // List of admins
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Message represents a message sent in a chat room.
type Message struct {
	ID        int       `json:"id"`
	RoomID    int       `json:"room_id"` // Associated room ID
	UserID    int       `json:"user_id"` // ID of the sender
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// RoomCreateRequest represents the payload for creating a new room.
type RoomCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

// RoomUpdateRequest represents the payload for updating room details.
type RoomUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Users       *[]int  `json:"users,omitempty"`
}

// RoomResponse represents the room object returned in API responses.
type RoomResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedBy   int       `json:"created_by"`
	Users       []int     `json:"users"`
	Messages    []Message `json:"messages"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MessageCreateRequest represents the payload for sending a message.
type MessageCreateRequest struct {
	RoomID  int    `json:"room_id" binding:"required"`
	UserID  int    `json:"user_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// MessageResponse represents the message object returned in API responses.
type MessageResponse struct {
	ID        int       `json:"id"`
	RoomID    int       `json:"room_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
