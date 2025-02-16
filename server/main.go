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
	// Connect to database
	db.ConnectDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	systemLogRepo := repository.NewSystemLogRepository(db.DB)
	roomRepo := repository.NewRoomRepository(db.DB)

	// Initialize services
	userService := services.NewUserService(userRepo)
	systemLogService := services.NewSystemLogService(systemLogRepo)
	roomService := services.NewRoomService(roomRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	systemLogHandler := handlers.NewSystemLogHandler(systemLogService)
	roomHandler := handlers.NewRoomHandler(roomService)
	wsHandler := handlers.NewWebSocketHandler(roomService) // WebSocket handler

	// Initialize router
	router := gin.Default()
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.SystemLogMiddleware()) // Middleware to log all requests

	// Setup routes (moved to app_routes.go)
	routes.SetupRoutes(router, userHandler, systemLogHandler, roomHandler, wsHandler)

	log.Println("🚀 Server started on port 8080")
	router.Run(":8080")
}
