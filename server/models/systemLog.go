package models

import "time"

// SystemLog represents a system log entry.
type SystemLog struct {
	ID         int        `json:"id"`
	Method     string     `json:"method"`
	Endpoint   string     `json:"endpoint"`
	UserID     *int       `json:"user_id,omitempty"`
	StatusCode int        `json:"status_code"`
	Message    string     `json:"message"`
	Timestamp  time.Time  `json:"timestamp"`
}