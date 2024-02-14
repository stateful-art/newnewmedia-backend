package productroutes

import (
	"github.com/gofiber/fiber/v2"
	controller "newnewmedia.com/microservices/place/controller"
)

func PlaceRoutes(app fiber.Router) {
	app.Get("/", controller.GetPlaces)
	app.Post("/", controller.CreatePlace)
	app.Get("/:id", controller.GetPlace)
}
