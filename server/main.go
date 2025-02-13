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

	userRepo := repository.NewUserRepository(db.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	router := gin.Default()
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.SystemLogMiddleware()) // Middleware to log all requests

	routes.SetupUserRoutes(router, userHandler)

	log.Println("🚀 Server started on port 8080")
	router.Run(":8080")
}