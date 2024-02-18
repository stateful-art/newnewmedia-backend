package productroutes

import (
	"github.com/gofiber/fiber/v2"
	controller "newnewmedia.com/microservices/music/controller"
)

func MusicRoutes(app fiber.Router) {
	app.Get("/", controller.GetMusic)
	app.Post("/", controller.CreateMusic)
	app.Get("/place/:id", controller.GetMusicByPlace)
	app.Get("/file/:id", controller.GetMusicFile)
	app.Get("/play/:id", controller.PlayMusic)
}
