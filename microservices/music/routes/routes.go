package musicroutes

import (
	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	controller "newnewmedia.com/microservices/music/controller"
)

func MusicRoutes(app fiber.Router, storageClient *storage.Client) {
	app.Post("/", func(c *fiber.Ctx) error {
		return controller.CreateMusic(c, storageClient)
	})

	app.Get("/place/:id", controller.GetMusicByPlace)

	app.Get("/play/:id", func(c *fiber.Ctx) error {
		return controller.PlayMusic(c, storageClient)
	})

	app.Get("/:id", func(c *fiber.Ctx) error {
		_, err := controller.GetSong(c, storageClient)
		return err
	})

}
