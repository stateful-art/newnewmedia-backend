package authroutes

import (
	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	controller "newnew.media/microservices/auth/controller"
)

func AuthRoutes(app fiber.Router, storageClient *storage.Client, redisClient *redis.Client) {

	// Register Spotify authentication routes
	app.Get("/spotify", controller.SpotifyLogin) // Initiate Spotify login

	// Handle Spotify callback after authorization
	app.Get("/spotify/callback", func(c *fiber.Ctx) error {
		return controller.SpotifyCallback(c, redisClient)
	})

}
