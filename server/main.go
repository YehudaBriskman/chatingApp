package main

import (
	"github.com/gin-gonic/gin"
	"github.com/YehudaBriskman/chatingApp/server/db"
	"github.com/YehudaBriskman/chatingApp/server/handlers"
	"github.com/YehudaBriskman/chatingApp/server/repository"
	"github.com/YehudaBriskman/chatingApp/server/routes"
	"github.com/YehudaBriskman/chatingApp/server/services"
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
