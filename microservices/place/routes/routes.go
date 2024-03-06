package placeroutes

import (
	"github.com/gofiber/fiber/v2"
	controller "newnew.media/microservices/place/controller"
)

func PlaceRoutes(app fiber.Router) {
	app.Get("/", controller.GetPlaces)
	app.Post("/", controller.CreatePlace)
	app.Get("/:id", controller.GetPlace)
}
