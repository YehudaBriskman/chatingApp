package repository

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"chatingApp/models"

	"github.com/lib/pq"
)

type RoomRepository struct {
	DB *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{DB: db}
}

// CreateRoom inserts a new chat room into the database and returns the created room
func (repo *RoomRepository) CreateRoom(room *models.Room) (*models.Room, error) {
	query := `INSERT INTO rooms (name, description, created_by, room_admins, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			  RETURNING id, name, description, created_by, room_admins, created_at, updated_at;`

	var roomAdminsStr string
	err := repo.DB.QueryRow(query, room.Name, room.Description, room.CreatedBy, pq.Array(room.RoomAdmins)).
		Scan(&room.ID, &room.Name, &room.Description, &room.CreatedBy, &roomAdminsStr, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// המרת room_admins ממחרוזת ל- []int
	room.RoomAdmins, err = parseIntArray(roomAdminsStr)
	if err != nil {
		return nil, err
	}

	return room, nil
}

// GetRoomByID retrieves a chat room by its ID
func (repo *RoomRepository) GetRoomByID(roomID int) (*models.Room, error) {
	query := `SELECT id, name, description, created_by, room_admins, created_at, updated_at FROM rooms WHERE id = $1;`
	row := repo.DB.QueryRow(query, roomID)

	var roomAdminsStr string
	room := &models.Room{}
	err := row.Scan(&room.ID, &room.Name, &room.Description, &room.CreatedBy, &roomAdminsStr, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// המרת room_admins ממחרוזת ל- []int
	room.RoomAdmins, err = parseIntArray(roomAdminsStr)
	if err != nil {
		return nil, err
	}

	return room, nil
}

// GetAllRooms retrieves a list of chat rooms from the database
func (repo *RoomRepository) GetAllRooms() ([]models.Room, error) {
	query := `SELECT id, name, description, created_by, room_admins, created_at, updated_at FROM rooms;`
	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		var roomAdminsStr string

		err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.CreatedBy, &roomAdminsStr, &room.CreatedAt, &room.UpdatedAt)
		if err != nil {
			return nil, err
		}

		room.RoomAdmins, err = parseIntArray(roomAdminsStr)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}
// DeleteRoom removes a chat room by its ID
func (repo *RoomRepository) DeleteRoom(roomID int) error {
	query := `DELETE FROM rooms WHERE id = $1;`
	_, err := repo.DB.Exec(query, roomID)
	return err
}

func (repo *RoomRepository) IsUserRoomAdmin(userID, roomID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM rooms WHERE id = $1 AND $2 = ANY(room_admins));`
	var exists bool
	err := repo.DB.QueryRow(query, roomID, userID).Scan(&exists)
	return exists, err
}

// UpdateRoomAdmins updates the list of admins for a room
func (repo *RoomRepository) UpdateRoomAdmins(roomID int, roomAdmins []int) error {
	query := `UPDATE rooms SET room_admins = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2;`
	_, err := repo.DB.Exec(query, pq.Array(roomAdmins), roomID)
	return err
}

// UpdateRoomDetails updates only the name and description of a room
func (repo *RoomRepository) UpdateRoomDetails(roomID int, name, description string) error {
	query := `UPDATE rooms SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3;`
	_, err := repo.DB.Exec(query, name, description, roomID)
	return err
}

// AddUserToRoom adds a user to a chat room
func (repo *RoomRepository) AddUserToRoom(roomID, userID int) error {
	query := `INSERT INTO room_users (room_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;`
	_, err := repo.DB.Exec(query, roomID, userID)
	return err
}

// RemoveUserFromRoom removes a user from a chat room
func (repo *RoomRepository) RemoveUserFromRoom(roomID, userID int) error {
	query := `DELETE FROM room_users WHERE room_id = $1 AND user_id = $2;`
	_, err := repo.DB.Exec(query, roomID, userID)
	return err
}

// IsUserInRoom checks if a user is a member of a chat room
func (repo *RoomRepository) IsUserInRoom(roomID, userID int) bool {
	query := `SELECT COUNT(*) FROM room_users WHERE room_id = $1 AND user_id = $2;`
	var count int
	err := repo.DB.QueryRow(query, roomID, userID).Scan(&count)
	return err == nil && count > 0
}

// GetUsersInRoom retrieves all users in a room
func (repo *RoomRepository) GetUsersInRoom(roomID int) ([]int, error) {
	query := `SELECT user_id FROM room_users WHERE room_id = $1;`
	rows, err := repo.DB.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		users = append(users, userID)
	}

	return users, nil
}

// parseIntArray converts a PostgreSQL array string "{1,2,3}" to a []int slice
func parseIntArray(pgArray string) ([]int, error) {
	pgArray = strings.Trim(pgArray, "{}")
	if pgArray == "" {
		return []int{}, nil
	}

	strValues := strings.Split(pgArray, ",")
	intValues := make([]int, len(strValues))

	for i, str := range strValues {
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		intValues[i] = val
	}

	return intValues, nil
}
