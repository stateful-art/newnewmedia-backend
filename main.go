package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
	utils "newnewmedia.com/commons/utils"
	db "newnewmedia.com/db"
	authroute "newnewmedia.com/microservices/auth/routes"
	musicroute "newnewmedia.com/microservices/music/routes"
	placesroute "newnewmedia.com/microservices/place/routes"
	playlistroute "newnewmedia.com/microservices/playlist/routes"
)

var StorageClient *storage.Client // Global variable to hold the GCS client instance
var RedisClient *redis.Client     // Global variable to hold the Redis client instance// Initialize the GCS client during application startup
func init() {

	if err := utils.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	_, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

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
	log.Println("Google Storage: OK")

	// Initialize Redis client with retry logic
	redisCtx := context.Background()
	// var redisOptions *redis.Options
	var redisErr error
	for attempt := 1; attempt <= 3; attempt++ { // Retry 3 times with exponential backoff
		RedisClient, _, redisErr = connectToRedis(redisCtx)
		if redisErr == nil {
			log.Println("Redis : OK")
			break
		}
		log.Printf("Failed to connect to Redis (attempt %d): %v\n", attempt, redisErr)
		time.Sleep(time.Duration(attempt) * time.Second) // Exponential backoff
	}
	if redisErr != nil {
		log.Fatalf("Failed to connect to Redis after multiple attempts: %v", redisErr)
	}

	// Optionally, you can log Redis client options
	// log.Printf("Redis Client Options: %+v\n", redisOptions)
}

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("No .env file found")
	// }
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "x-spotify-token ,Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		ExposeHeaders:    "Content-Length,Access-Control-Allow-Headers", // Expose the required header
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
	playlistroute.PlaylistRoutes(playlists)

	log.Print("APP @ 3000 : OK")
	log.Fatal(app.Listen(":3000"))
}

func connectToRedis(ctx context.Context) (*redis.Client, *redis.Options, error) {
	options := &redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"), // Redis address
		Password: "",                         // Redis password, if any
		DB:       0,                          // Redis database index
	}
	client := redis.NewClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, options, err
	}
	return client, options, nil
}
