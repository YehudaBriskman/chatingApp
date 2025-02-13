package repository

import (
	"database/sql"
	"log"
	"time"
	"chatingApp/models"
)

// SystemLogRepository handles database operations for system logs.
type SystemLogRepository struct {
	DB *sql.DB
}

// NewSystemLogRepository initializes a new SystemLogRepository instance.
func NewSystemLogRepository(db *sql.DB) *SystemLogRepository {
	return &SystemLogRepository{DB: db}
}

// AddLog inserts a new system log entry into the database.
func (repo *SystemLogRepository) AddLog(method, endpoint string, userID *int, statusCode int, message string) error {
	query := "INSERT INTO system_logs (method, endpoint, user_id, status_code, message, timestamp) VALUES ($1, $2, $3, $4, $5, $6)"
	timestamp := time.Now()

	_, err := repo.DB.Exec(query, method, endpoint, userID, statusCode, message, timestamp)
	if err != nil {
		log.Println("Error: Failed to insert system log", err)
		return err
	}

	log.Println("✅ System log added successfully:", method, endpoint, statusCode, message)
	return nil
}

// GetAllLogs retrieves all system logs from the database.
func (repo *SystemLogRepository) GetAllLogs() ([]models.SystemLog, error) {
	rows, err := repo.DB.Query("SELECT id, method, endpoint, user_id, status_code, message, timestamp FROM system_logs")
	if err != nil {
		log.Println("Error: Failed to retrieve system logs", err)
		return nil, err
	}
	defer rows.Close()

	var logs []models.SystemLog

	for rows.Next() {
		var logEntry models.SystemLog
		if err := rows.Scan(&logEntry.ID, &logEntry.Method, &logEntry.Endpoint, &logEntry.UserID, &logEntry.StatusCode, &logEntry.Message, &logEntry.Timestamp); err != nil {
			return nil, err
		}
		logs = append(logs, logEntry)
	}

	return logs, nil
}

// GetLogsByUser retrieves logs associated with a specific user.
func (repo *SystemLogRepository) GetLogsByUser(userID int) ([]models.SystemLog, error) {
	rows, err := repo.DB.Query("SELECT id, method, endpoint, user_id, status_code, message, timestamp FROM system_logs WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Error: Failed to retrieve system logs for user", err)
		return nil, err
	}
	defer rows.Close()

	var logs []models.SystemLog

	for rows.Next() {
		var logEntry models.SystemLog
		if err := rows.Scan(&logEntry.ID, &logEntry.Method, &logEntry.Endpoint, &logEntry.UserID, &logEntry.StatusCode, &logEntry.Message, &logEntry.Timestamp); err != nil {
			return nil, err
		}
		logs = append(logs, logEntry)
	}

	return logs, nil
}