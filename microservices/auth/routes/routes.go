package authroutes

import (
	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	controller "newnewmedia.com/microservices/auth/controller"
)

func AuthRoutes(app fiber.Router, storageClient *storage.Client) {

	// Register Spotify authentication routes
	app.Get("/spotify", controller.SpotifyLogin)             // Initiate Spotify login
	app.Get("/spotify/callback", controller.SpotifyCallback) // Handle Spotify callback after authorization
}
