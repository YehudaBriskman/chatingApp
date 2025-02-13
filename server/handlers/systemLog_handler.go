package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"chatingApp/services"
	"chatingApp/middleware"
	"github.com/gin-gonic/gin"
)

// LogHandler handles HTTP requests for system logs.
type LogHandler struct {
	LogService *services.SystemLogService
}

// NewLogHandler creates a new LogHandler instance.
func NewSystemLogHandler(service *services.SystemLogService) *LogHandler {
	return &LogHandler{LogService: service}
}

// GetLogs handles the GET request to retrieve all system logs.
func (h *LogHandler) GetLogs(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication token"})
		return
	}

	claims, err := services.ValidateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	if !middleware.HasRequiredRole(claims["role"].(string), "admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	logs, err := h.LogService.GetAllLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// GetLogsByUser handles the GET request to retrieve logs by user ID.
func (h *LogHandler) GetLogsByUser(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication token"})
		return
	}

	claims, err := services.ValidateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	if !middleware.HasRequiredRole(claims["role"].(string), "admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	userIDParam := c.Param("userID")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	logs, err := h.LogService.GetLogsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs for user"})
		return
	}
	c.JSON(http.StatusOK, logs)
}
