package main

import (
	"chatingApp/db"
	"chatingApp/handlers"
	"chatingApp/middleware"
	"chatingApp/repository"
	"chatingApp/routes"
	"chatingApp/services"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	db.ConnectDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	systemLogRepo := repository.NewSystemLogRepository(db.DB)

	// Initialize services
	userService := services.NewUserService(userRepo)
	systemLogService := services.NewSystemLogService(systemLogRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	systemLogHandler := handlers.NewSystemLogHandler(systemLogService)

	// Initialize router
	router := gin.Default()
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.SystemLogMiddleware()) // Middleware to log all requests

	// Setup routes
	routes.SetupUserRoutes(router, userHandler)
	routes.SetupLogRoutes(router, systemLogHandler)

	log.Println("🚀 Server started on port 8080")
	router.Run(":8080")
}
