package main

import (
	"github.com/gin-gonic/gin"
	"chatingApp/db"
	"chatingApp/handlers"
	"chatingApp/repository"
	"chatingApp/routes"
	"chatingApp/services"
)

func main() {
	db.ConnectDB()

	userRepo := repository.NewUserRepository(db.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	router := gin.Default()
	routes.SetupUserRoutes(router, userHandler)

	router.Run(":8080")
}
