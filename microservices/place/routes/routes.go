package placeroutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	controller "newnew.media/microservices/place/controller"
	repo "newnew.media/microservices/place/repository"
	service "newnew.media/microservices/place/service"
)

func PlaceRoutes(app fiber.Router, natsClient *nats.Conn) {

	repo := repo.NewPlaceRepository()
	// Create an instance of PlaylistService with the PlaylistRepository
	service := service.NewPlaceService(repo, natsClient)

	// Pass the PlaylistService instance to the PlaylistController
	controller := controller.NewPlaceController(service)

	app.Get("/:id", controller.GetPlace)
	app.Get("/", controller.GetPlaces)
	app.Get("/nearby/:lat/:long", controller.GetPlacesNearLocation)
	app.Post("/", controller.CreatePlace)

	// utility endpoints
	app.Get("/index/create-geo", controller.CreateGeospatialIndex)

}
