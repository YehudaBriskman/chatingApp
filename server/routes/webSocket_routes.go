package routes

import (
	"chatingApp/handlers"
	"chatingApp/middleware"
	"github.com/gin-gonic/gin"
)

// SetupWebSocketRoutes configures WebSocket routes for real-time chat.
func SetupWebSocketRoutes(router *gin.Engine, wsHandler *handlers.WebSocketHandler) {
	wsRoutes := router.Group("/ws")
	{
		wsRoutes.GET("/:roomID", middleware.AuthMiddleware(), wsHandler.HandleWebSocketConnection)
	}
}
