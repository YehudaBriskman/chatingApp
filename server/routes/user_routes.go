package routes

import (
	"chatingApp/handlers"
	"chatingApp/middleware"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", middleware.AuthMiddleware(), middleware.AdminMiddleware("admin"), userHandler.GetUsers) // Only admin or higher can access
		userRoutes.POST("/add", userHandler.AddUser) // Requires authentication
		userRoutes.POST("/login", userHandler.Login) // Open for all
	}
}