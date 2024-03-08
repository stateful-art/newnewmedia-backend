package userroutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	controller "newnew.media/microservices/user/controller"
	repo "newnew.media/microservices/user/repository"
	service "newnew.media/microservices/user/service"
)

func UserRoutes(app fiber.Router, redisClient *redis.Client, natsClient *nats.Conn) {
	// Create an instance of PlaylistRepository
	userRepo := repo.NewUserRepository()
	// smtpService := communicationService.NewSMTPService(redisClient)
	// Create an instance of PlaylistService with the PlaylistRepository
	userService := service.NewUserService(userRepo, redisClient, natsClient)

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
