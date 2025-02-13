package main

import (
	"chatingApp/db"
	"chatingApp/handlers"
	"chatingApp/middleware"
	"chatingApp/repository"
	"chatingApp/routes"
	"chatingApp/services"
	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDB()

	userRepo := repository.NewUserRepository(db.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	router := gin.Default()
	router.Use(middleware.ErrorHandlerMiddleware())
	routes.SetupUserRoutes(router, userHandler)

	router.Run(":8080")
}
