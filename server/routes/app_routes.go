package routes

import (
	"chatingApp/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, userHandler *handlers.UserHandler, logHandler *handlers.LogHandler, roomHandler *handlers.RoomHandler, wsHandler *handlers.WebSocketHandler) {
	// User & Log Routes
	SetupUserRoutes(router, userHandler)
	SetupLogRoutes(router, logHandler)

	// Room & WebSocket Routes
	SetupRoomRoutes(router, roomHandler)
	SetupWebSocketRoutes(router, wsHandler)
}
