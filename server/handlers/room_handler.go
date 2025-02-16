package handlers

import (
	"chatingApp/middleware"
	"chatingApp/models"
	"chatingApp/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// RoomHandler handles HTTP requests for room operations.
type RoomHandler struct {
	RoomService *services.RoomService
}

// NewRoomHandler creates a new RoomHandler instance.
func NewRoomHandler(service *services.RoomService) *RoomHandler {
	return &RoomHandler{RoomService: service}
}

// CreateRoom handles the POST request to create a new room.
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	userID, _, _, err := middleware.ExtractTokenData(c, "admin")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var roomInput models.RoomCreateRequest
	if err := c.ShouldBindJSON(&roomInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	room := models.Room{
		Name:        roomInput.Name,
		Description: roomInput.Description,
		CreatedBy:   userID,
		RoomAdmins:  []int{userID},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdRoom, err := h.RoomService.CreateRoom(&room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}

	err = h.RoomService.AddUserToRoom(createdRoom.ID, userID, userID)
	if err != nil {
		delErr := h.RoomService.DeleteRoom(createdRoom.ID, userID)
		if delErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":       "Failed to update room users status",
				"stack":       err.Error(),
				"deleteError": delErr.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update room users status",
				"stack": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Room created successfully", "room": createdRoom})
}

// GetRooms handles the GET request to retrieve all rooms.
func (h *RoomHandler) GetRooms(c *gin.Context) {
	rooms, err := h.RoomService.GetAllRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms"})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

// GetRoom handles the GET request to retrieve a room by ID.
func (h *RoomHandler) GetRoom(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	room, err := h.RoomService.GetRoom(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch room"})
		return
	}
	c.JSON(http.StatusOK, room)
}

// DeleteRoom handles the DELETE request to remove a room.
func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	userID, _, _, err := middleware.ExtractTokenData(c, "admin")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	err = h.RoomService.DeleteRoom(roomID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
}

func (h *RoomHandler) IsUserRoomAdmin(c *gin.Context) {
	userID, _, _, err := middleware.ExtractTokenData(c, "user")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}
	exist,err := h.RoomService.IsUserRoomAdmin(roomID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check admin status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"exist": exist})
}


// // UpdateRoomDetails handles the PUT request to update room details.
// func (h *RoomHandler) UpdateRoomDetails(c *gin.Context) {
// 	var updateInput models.RoomUpdateRequest

// 	roomID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
// 		return
// 	}

// 	userID := c.GetInt("user_id")
// 	if !h.RoomService.IsUserRoomAdmin(roomID, userID) {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
// 		return
// 	}

// 	if err := c.ShouldBindJSON(&updateInput); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	err = h.RoomService.UpdateRoomDetails(roomID, *updateInput.Name, *updateInput.Description, userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room details"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Room details updated successfully"})
// }

// // UpdateRoomAdmins handles the PUT request to update room admins.
// func (h *RoomHandler) UpdateRoomAdmins(c *gin.Context) {
// 	var adminInput struct {
// 		RoomAdmins []int `json:"room_admins" binding:"required"`
// 	}

// 	roomID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
// 		return
// 	}

// 	userID := c.GetInt("user_id")
// 	if !h.RoomService.IsUserRoomAdmin(roomID, userID) {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
// 		return
// 	}

// 	if err := c.ShouldBindJSON(&adminInput); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	err = h.RoomService.UpdateRoomAdmins(roomID, adminInput.RoomAdmins, userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room admins"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Room admins updated successfully"})
// }

// // AddUserToRoom handles adding a user to a room.
// func (h *RoomHandler) AddUserToRoom(c *gin.Context) {
// 	var input struct {
// 		UserID int `json:"user_id" binding:"required"`
// 	}

// 	roomID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
// 		return
// 	}

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	err = h.RoomService.AddUserToRoom(roomID, input.UserID, c.GetInt("user_id"))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to room"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "User added successfully to room"})
// }

// // GetUsersInRoom handles the GET request to retrieve all users in a room.
// func (h *RoomHandler) GetUsersInRoom(c *gin.Context) {
// 	roomID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
// 		return
// 	}

// 	users, err := h.RoomService.GetUsersInRoom(roomID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users in room"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"users": users})
// }

// // GetRoomsByUserID handles the GET request to retrieve all rooms a user is a member of.
// func (h *RoomHandler) GetRoomsByUserID(c *gin.Context) {
// 	userID, err := strconv.Atoi(c.Param("user_id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
// 		return
// 	}

// 	rooms, err := h.RoomService.GetRoomsByUserID(userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rooms"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
// }

// // GetMessagesByRoomID handles the GET request to retrieve all messages in a room.
// func (h *RoomHandler) GetMessagesByRoomID(c *gin.Context) {
// 	roomID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
// 		return
// 	}

// 	messages, err := h.RoomService.GetMessagesByRoomID(roomID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"messages": messages})
// }

// // DeleteMessageByID handles the DELETE request to remove a message from a chat room.
// func (h *RoomHandler) DeleteMessageByID(c *gin.Context) {
// 	messageID, err := strconv.Atoi(c.Param("message_id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
// 		return
// 	}

// 	userID := c.GetInt("user_id")
// 	err = h.RoomService.DeleteMessageByID(messageID, userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
// }

// // IsUserInRoom handles the GET request to check if a user is in a room.
// func (h *RoomHandler) IsUserInRoom(c *gin.Context) {
// 	roomID, err := strconv.Atoi(c.Param("room_id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
// 		return
// 	}

// 	userID, err := strconv.Atoi(c.Param("user_id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
// 		return
// 	}

// 	isInRoom := h.RoomService.IsUserInRoom(roomID, userID)
// 	c.JSON(http.StatusOK, gin.H{"room_id": roomID, "user_id": userID, "is_in_room": isInRoom})
// }
