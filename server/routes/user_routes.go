package routes

import (
	"chatingApp/handlers"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", userHandler.GetUsers)
		userRoutes.POST("/add", userHandler.AddUser)
		userRoutes.POST("/login", userHandler.Login)
	}
}
