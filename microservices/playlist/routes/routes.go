package playlistroutes

import (
	"github.com/gofiber/fiber/v2"
	controller "newnewmedia.com/microservices/playlist/controller"
	repo "newnewmedia.com/microservices/playlist/repository"
	service "newnewmedia.com/microservices/playlist/service"
)

func PlaylistRoutes(app fiber.Router) {
	// Create an instance of PlaylistRepository
	playlistRepo := repo.NewPlaylistRepository()

	// Create an instance of PlaylistService with the PlaylistRepository
	playlistService := service.NewPlaylistService(playlistRepo)

	// Pass the PlaylistService instance to the PlaylistController
	playlistController := controller.NewPlaylistController(playlistService)

	// Define routes using the controller methods
	app.Get("/", playlistController.GetPlaylists)
	app.Get("/:id", playlistController.GetPlaylistByID)

	app.Post("/", playlistController.CreatePlaylist)
	app.Patch("/:id", playlistController.UpdatePlaylist)
	app.Delete("/:id", playlistController.DeletePlaylist)

}
