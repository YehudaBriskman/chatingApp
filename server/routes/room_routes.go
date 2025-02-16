package routes

import (
	"chatingApp/handlers"
	"chatingApp/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoomRoutes configures routes for chat room management.
func SetupRoomRoutes(router *gin.Engine, roomHandler *handlers.RoomHandler) {
	roomRoutes := router.Group("/rooms")
	{
		roomRoutes.POST("/", middleware.AuthMiddleware(), roomHandler.CreateRoom)
		roomRoutes.GET("/", middleware.AuthMiddleware(), middleware.AdminMiddleware("super-admin"), roomHandler.GetRooms)
		roomRoutes.GET("/:id", middleware.AuthMiddleware(), roomHandler.GetRoom)
		roomRoutes.DELETE("/:id", middleware.AuthMiddleware(), roomHandler.DeleteRoom)
		roomRoutes.GET("/room/:id", middleware.AuthMiddleware(), roomHandler.IsUserRoomAdmin)
		// roomRoutes.PUT("/:id", middleware.AuthMiddleware(), roomHandler.UpdateRoomDetails)
		// roomRoutes.PUT("/:id/admins", middleware.AuthMiddleware(), roomHandler.UpdateRoomAdmins)
		// roomRoutes.POST("/:id/users", middleware.AuthMiddleware(), roomHandler.AddUserToRoom)
	}
}
