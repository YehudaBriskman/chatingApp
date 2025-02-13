package services

import (
	"chatingApp/models"
	"chatingApp/repository"
	"log"
)

// SystemLogService provides business logic for system logs.
type SystemLogService struct {
	LogRepo *repository.SystemLogRepository
}

// NewSystemLogService creates a new instance of SystemLogService.
func NewSystemLogService(repo *repository.SystemLogRepository) *SystemLogService {
	return &SystemLogService{LogRepo: repo}
}

// AddLog adds a new system log entry.
func (s *SystemLogService) AddLog(method, endpoint string, userID *int, statusCode int, message string) error {
	err := s.LogRepo.AddLog(method, endpoint, userID, statusCode, message)
	if err != nil {
		log.Println("❌ Error: Failed to log system event", err)
		return err
	}
	log.Println("✅ System log recorded successfully:", method, endpoint, statusCode, message)
	return nil
}

// GetAllLogs retrieves all system logs.
func (s *SystemLogService) GetAllLogs() ([]models.SystemLog, error) {
	logs, err := s.LogRepo.GetAllLogs()
	if err != nil {
		log.Println("❌ Error: Failed to retrieve system logs", err)
		return nil, err
	}
	return logs, nil
}

// GetLogsByUser retrieves system logs by a specific user ID.
func (s *SystemLogService) GetLogsByUser(userID int) ([]models.SystemLog, error) {
	logs, err := s.LogRepo.GetLogsByUser(userID)
	if err != nil {
		log.Println("❌ Error: Failed to retrieve system logs for user", err)
		return nil, err
	}
	return logs, nil
}
