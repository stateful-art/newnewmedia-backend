package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "newnewmedia.com/microservices/playlist/controller"
)

func PlaceRoutes(app fiber.Router) {
	app.Get("/", controller.GetPlaylists)
}
