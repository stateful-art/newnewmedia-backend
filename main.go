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
	"github.com/redis/go-redis/v9"
	db "newnewmedia.com/db"
	authroute "newnewmedia.com/microservices/auth/routes"
	musicroute "newnewmedia.com/microservices/music/routes"
	placesroute "newnewmedia.com/microservices/place/routes"
	playlistroute "newnewmedia.com/microservices/playlist/routes"
)

var StorageClient *storage.Client // Global variable to hold the GCS client instance
var RedisClient *redis.Client     // Global variable to hold the Redis client instance// Initialize the GCS client during application startup
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

	redisNodes := os.Getenv("REDIS_CLUSTER_NODES")
	if redisNodes == "" {
		log.Fatal("REDIS_CLUSTER_NODES environment variable is not set")
	}

	// Initialize Redis client
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"), // Redis address
		Password: "",                         // Redis password, if any
		DB:       0,                          // Redis database index
	})
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	app := fiber.New()
	db.ConnectDB()

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return origin == "https://www.app.newnew.media/" || origin == "http://localhost:5173/"
		},
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "x-spotify-token ,Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowCredentials: true,
	}))
	app.Use(logger.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	auth := app.Group("/auth")
	music := app.Group("/music")
	places := app.Group("/places")
	playlists := app.Group("/playlists")

	authroute.AuthRoutes(auth, StorageClient, RedisClient)
	placesroute.PlaceRoutes(places)
	musicroute.MusicRoutes(music, StorageClient)
	playlistroute.PlaceRoutes(playlists)

	log.Fatal(app.Listen(":3000"))
}
