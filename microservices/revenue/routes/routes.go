package revenueroutes

import (
	"github.com/gofiber/fiber/v2"
	controller "newnewmedia.com/microservices/revenue/controller"
	repo "newnewmedia.com/microservices/revenue/repository"
	service "newnewmedia.com/microservices/revenue/service"
)

func RevenueRoutes(app fiber.Router) {
	// Create an instance of PlaylistRepository
	revenueRepo := repo.NewRevenueRepository()

	// Create an instance of PlaylistService with the PlaylistRepository
	revenueService := service.NewRevenueService(revenueRepo)

	// Pass the PlaylistService instance to the revenueController
	revenueController := controller.NewRevenueController(revenueService)

	app.Get("/:id", revenueController.GetRevenueByID)
	app.Get("/:artist_id", revenueController.GetRevenueByArtistID)
	app.Get("/:playlist_id", revenueController.GetRevenueByPlaylistID)
	app.Get("/:place_id", revenueController.GetRevenueByPlaceID)

	app.Post("/", revenueController.CreateRevenue)

}
