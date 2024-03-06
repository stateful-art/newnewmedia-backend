package playlistroutes

import (
	"github.com/gofiber/fiber/v2"
	controller "newnew.media/microservices/playlist/controller"
	repo "newnew.media/microservices/playlist/repository"
	service "newnew.media/microservices/playlist/service"
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
	app.Get("/:owner_id", playlistController.GetPlaylistsByOwner)

	app.Post("/", playlistController.CreatePlaylist)
	app.Patch("/:id", playlistController.UpdatePlaylist)
	app.Patch("/add-songs", playlistController.AddSongsToPlaylist)
	app.Patch("/remove-songs", playlistController.RemoveSongsFromPlaylist)
	app.Delete("/:id", playlistController.DeletePlaylist)

}
