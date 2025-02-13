package middleware

import (
	"log"
	"time"
	"chatingApp/db"
	"github.com/gin-gonic/gin"
)

// SystemLogMiddleware logs all incoming requests to the database
func SystemLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		method := c.Request.Method
		endpoint := c.Request.URL.Path
		statusCode := c.Writer.Status()
		message := "Request processed"

		userID, exists := c.Get("userID")
		var userIDValue interface{}
		if exists {
			userIDValue = userID
		} else {
			userIDValue = nil
		}

		_, err := db.DB.Exec(`
			INSERT INTO system_logs (method, endpoint, user_id, status_code, message, timestamp) 
			VALUES ($1, $2, $3, $4, $5, $6)`,
			method, endpoint, userIDValue, statusCode, message, startTime,
		)
		if err != nil {
			log.Printf("❌ Error: failed to insert log into system_logs: %v", err)
		}
	}
}
