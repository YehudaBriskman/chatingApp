package services

import (
	"chatingApp/models"
	"chatingApp/repository"
	"errors"

	// "errors"
	"log"
)

// RoomService provides business logic for chat rooms.
type RoomService struct {
	RoomRepo *repository.RoomRepository
}

// NewRoomService creates a new instance of RoomService.
func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{RoomRepo: repo}
}

// CreateRoom creates a new chat room.
func (s *RoomService) CreateRoom(room *models.Room) (*models.Room, error) {
	createdRoom, err := s.RoomRepo.CreateRoom(room)
	if err != nil {
		log.Println("❌ Error: Failed to create room", err)
		return nil, err
	}
	log.Println("✅ Room created successfully:", createdRoom.Name)
	return createdRoom, nil
}

// GetRoom retrieves a room by ID.
func (s *RoomService) GetRoom(roomID int) (*models.Room, error) {
	room, err := s.RoomRepo.GetRoomByID(roomID)
	if err != nil {
		log.Println("❌ Error: Failed to retrieve room", err)
		return nil, err
	}
	return room, nil
}

// GetAllRooms retrieves all rooms.
func (s *RoomService) GetAllRooms() ([]models.Room, error) {
	rooms, err := s.RoomRepo.GetAllRooms()
	if err != nil {
		log.Println("❌ Error: Failed to retrieve rooms", err)
		return nil, err
	}
	return rooms, nil
}

// DeleteRoom removes a chat room.
func (s *RoomService) DeleteRoom(roomID, requesterID int) error {
	// Only room admins can delete rooms
	exist, err := s.IsUserRoomAdmin(roomID, requesterID)
	if (err != nil || !exist) {
		log.Println("❌ Error: User is not an admin of the room")
		return errors.New("user is not an admin of the room")
	}

	err = s.RoomRepo.DeleteRoom(roomID)
	if err != nil {
		log.Println("❌ Error: Failed to delete room", err)
		return err
	}
	log.Println("✅ Room deleted successfully:", roomID)
	return nil
}

// IsUserRoomAdmin checks if a user is an admin of a room.
func (s *RoomService) IsUserRoomAdmin(roomID, userID int) (bool, error) {
	exist, err := s.RoomRepo.IsUserRoomAdmin(userID, roomID)
	if err != nil {
		log.Println("❌ Error: Failed to check admin status", err)
		return false, err
	}
	return exist, nil
}


// // IsUserInRoom checks if a user is a member of a room.
// func (s *RoomService) IsUserInRoom(roomID, userID int) bool {
// 	return s.RoomRepo.IsUserInRoom(roomID, userID)
// }

// // UpdateRoomDetails updates name and description of a room.
// func (s *RoomService) UpdateRoomDetails(roomID int, name, description string, userID int) error {
// 	// Only room admins can update room details
// 	if !s.IsUserRoomAdmin(roomID, userID) {
// 		log.Println("❌ Error: User is not an admin of the room")
// 		return errors.New("user is not an admin of the room")
// 	}

// 	err := s.RoomRepo.UpdateRoomDetails(roomID, name, description)
// 	if err != nil {
// 		log.Println("❌ Error: Failed to update room details", err)
// 		return err
// 	}
// 	log.Println("✅ Room details updated successfully:", roomID)
// 	return nil
// }

// // UpdateRoomAdmins updates the admin list of a room.
// func (s *RoomService) UpdateRoomAdmins(roomID int, roomAdmins []int, userID int) error {
// 	// Only current admins can update the admin list
// 	if !s.IsUserRoomAdmin(roomID, userID) {
// 		log.Println("❌ Error: User is not an admin of the room")
// 		return errors.New("user is not an admin of the room")
// 	}

// 	err := s.RoomRepo.UpdateRoomAdmins(roomID, roomAdmins)
// 	if err != nil {
// 		log.Println("❌ Error: Failed to update room admins", err)
// 		return err
// 	}
// 	log.Println("✅ Room admins updated successfully:", roomID)
// 	return nil
// }

// AddUserToRoom adds a user to a chat room.
func (s *RoomService) AddUserToRoom(roomID, userID, requesterID int) error {
	// Only room admins can add users
	// if !s.IsUserRoomAdmin(roomID, requesterID) {
	// 	log.Println("❌ Error: User is not an admin of the room")
	// 	return errors.New("user is not an admin of the room")
	// }

	err := s.RoomRepo.AddUserToRoom(roomID, userID)
	if err != nil {
		log.Println("❌ Error: Failed to add user to room", err)
		return err
	}
	log.Println("✅ User added successfully to room:", roomID)
	return nil
}

// // RemoveUserFromRoom removes a user from a chat room.
// func (s *RoomService) RemoveUserFromRoom(roomID, userID, requesterID int) error {
// 	// Only room admins can remove users
// 	if !s.IsUserRoomAdmin(roomID, requesterID) {
// 		log.Println("❌ Error: User is not an admin of the room")
// 		return errors.New("user is not an admin of the room")
// 	}

// 	err := s.RoomRepo.RemoveUserFromRoom(roomID, userID)
// 	if err != nil {
// 		log.Println("❌ Error: Failed to remove user from room", err)
// 		return err
// 	}
// 	log.Println("✅ User removed successfully from room:", roomID)
// 	return nil
// }

// // AddMessageToRoom adds a message to a chat room.
// func (s *RoomService) AddMessageToRoom(roomID, userID int, content string) error {
// 	// Ensure user is in the room before adding a message
// 	if !s.RoomRepo.IsUserInRoom(roomID, userID) {
// 		log.Println("❌ Error: User is not in the room")
// 		return errors.New("user is not in the room")
// 	}

// 	err := s.RoomRepo.AddMessageToRoom(roomID, userID, content)
// 	if err != nil {
// 		log.Println("❌ Error: Failed to add message to room", err)
// 		return err
// 	}
// 	log.Println("✅ Message added successfully to room:", roomID)
// 	return nil
// }

// // GetUsersInRoom retrieves all users in a room.
// func (s *RoomService) GetUsersInRoom(roomID int) ([]int, error) {
// 	users, err := s.RoomRepo.GetUsersInRoom(roomID)
// 	if err != nil {
// 		log.Println("❌ Error: Failed to retrieve users in room", err)
// 		return nil, err
// 	}
// 	return users, nil
// }

// // GetRoomsByUserID retrieves all rooms a user is a member of.
// func (s *RoomService) GetRoomsByUserID(userID int) ([]models.Room, error) {
// 	rooms, err := s.RoomRepo.GetRoomsByUserID(userID)
// 	if err != nil {
// 		log.Println("❌ Error: Failed to retrieve rooms for user", err)
// 		return nil, err
// 	}
// 	return rooms, nil
// }

// // GetMessagesByRoomID retrieves all messages from a chat room.
// func (s *RoomService) GetMessagesByRoomID(roomID int) ([]models.Message, error) {
// 	messages, err := s.RoomRepo.GetMessagesByRoomID(roomID)
// 	if err != nil {
// 		log.Println("❌ Error: Failed to retrieve messages for room", err)
// 		return nil, err
// 	}
// 	return messages, nil
// }

// // DeleteMessageByID deletes a message from a chat room.
// func (s *RoomService) DeleteMessageByID(messageID, requesterID int) error {
// 	err := s.RoomRepo.DeleteMessageByID(messageID)
// 	if err != nil {
// 		log.Println("❌ Error: Failed to delete message", err)
// 		return err
// 	}
// 	log.Println("✅ Message deleted successfully:", messageID)
// 	return nil
// }
