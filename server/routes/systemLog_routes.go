package routes

import (
	"chatingApp/handlers"
	"chatingApp/middleware"
	"github.com/gin-gonic/gin"
)

func SetupLogRoutes(router *gin.Engine, logHandler *handlers.LogHandler) {
	logRoutes := router.Group("/logs")
	{
		logRoutes.GET("/", middleware.AuthMiddleware(), middleware.AdminMiddleware("admin"), logHandler.GetLogs) // Only admin or higher can access
		logRoutes.GET("/user/:userID", middleware.AuthMiddleware(), middleware.AdminMiddleware("admin"), logHandler.GetLogsByUser) // Only admin can access logs by user
	}
}