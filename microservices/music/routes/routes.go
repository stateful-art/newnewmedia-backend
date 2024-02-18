package productroutes

import (
	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	controller "newnewmedia.com/microservices/music/controller"
)

func MusicRoutes(app fiber.Router, storageClient *storage.Client) {
	app.Get("/", controller.GetMusic)
	app.Post("/", func(c *fiber.Ctx) error {
		return controller.CreateMusic(c, storageClient)
	})

	app.Get("/place/:id", controller.GetMusicByPlace)
	app.Get("/file/:id", controller.GetMusicFile)

	app.Get("/play/:id", func(c *fiber.Ctx) error {
		return controller.PlayMusic(c, storageClient)
	})
}
