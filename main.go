package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	db "newnewmedia.com/db"
	musicroute "newnewmedia.com/microservices/music/routes"
	placesroute "newnewmedia.com/microservices/place/routes"
	playlistroute "newnewmedia.com/microservices/playlist/routes"
)

var StorageClient *storage.Client // Global variable to hold the GCS client instance

// Initialize the GCS client during application startup
func init() {
	// Set the path to your credentials file
	credentialsFile := filepath.FromSlash("creds/creds.json")

	// Set the GOOGLE_APPLICATION_CREDENTIALS environment variable
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credentialsFile)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Println("Failed to initialize GCS client:", err)
		os.Exit(1)
	}
	StorageClient = client
	log.Println("Initialised google cloud storage client")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	app := fiber.New()
	db.ConnectDB()

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return origin == "https://www.newnewmedia.com" || origin == "http://localhost:5173"
		},
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowCredentials: true,
	}))
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	music := app.Group("/music")
	places := app.Group("/places")
	playlists := app.Group("/playlists")

	placesroute.PlaceRoutes(places)
	musicroute.MusicRoutes(music, StorageClient)
	playlistroute.PlaceRoutes(playlists)

	log.Fatal(app.Listen(":3000"))
}
