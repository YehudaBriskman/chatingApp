package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"os"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() 

		if len(c.Errors) > 0 {
			err := c.Errors.Last() 
			errorResponse := gin.H{"error": err.Error()}

			if os.Getenv("MODE") == "production" {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "there was an error"})
				return
			}

			c.JSON(http.StatusInternalServerError, errorResponse)
		}
	}
}
