package productroutes

import (
	"github.com/gofiber/fiber/v2"
	controller "newnewmedia.com/microservices/music/controller"
)

func MusicRoutes(app fiber.Router) {
	app.Get("/", controller.GetMusic)
}
