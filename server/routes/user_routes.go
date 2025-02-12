package routes

import (
	"github.com/YehudaBriskman/chatingApp/server/handlers"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", userHandler.GetUsers)
	}
}
