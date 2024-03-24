package musicroutes

import (
	"log"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	controller "newnew.media/microservices/music/controller"
)

func MusicRoutes(app fiber.Router, storageClient *storage.Client, redisClient *redis.Client) {
	app.Post("/", func(c *fiber.Ctx) error {
		log.Print("@ / route, about to call controller.CreateMusic(c, storageClient) ")
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

	// SPOTIFY RELATED ROUTES
	app.Get("/spotify/playlists", controller.UserPlaylists)
	app.Get("/spotify/recently-played-songs", controller.RecentlyPlayedSongs)
	app.Get("/spotify/genre-analysis", controller.GenreAnalysis)

}
