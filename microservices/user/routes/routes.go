package userroutes

import (
	"github.com/gofiber/fiber/v2"
	controller "newnew.media/microservices/user/controller"
	repo "newnew.media/microservices/user/repository"
	service "newnew.media/microservices/user/service"
)

func UserRoutes(app fiber.Router) {
	// Create an instance of PlaylistRepository
	userRepo := repo.NewUserRepository()

	// Create an instance of PlaylistService with the PlaylistRepository
	userService := service.NewUserService(userRepo)

	// Pass the PlaylistService instance to the PlaylistController
	userController := controller.NewUserController(userService)

	// Define routes using the controller methods
	app.Get("/", userController.GetUsers)
	app.Get("/:id", userController.GetUserByID)
	app.Get("/:spotify_id", userController.GetUserBySpotifyID)
	app.Get("/:youtube_id", userController.GetUserByYouTubeID)

	app.Post("/", userController.CreateUser)
	app.Patch("/add-role", userController.AddRole)
	app.Patch("/remove-role", userController.RemoveRole)

}
