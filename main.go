package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	utils "newnew.media/commons/utils"
	db "newnew.media/db"
	authroute "newnew.media/microservices/auth/routes"
	communicationroutes "newnew.media/microservices/communication/routes"
	musicroute "newnew.media/microservices/music/routes"
	placesroute "newnew.media/microservices/place/routes"
	playlistroute "newnew.media/microservices/playlist/routes"
	revenueroute "newnew.media/microservices/revenue/routes"
	searchroute "newnew.media/microservices/search/routes"
	userroute "newnew.media/microservices/user/routes"
)

var StorageClient *storage.Client
var RedisClient *redis.Client
var NatsClient *nats.Conn
var ElasticClient *elasticsearch.Client

func init() {
	if err := utils.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// log.Print("SENDING EMAIL w MAILGUN...")

	// commservice.SendEmail()

	_, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Initialize Google Cloud Storage client
	initGoogleCloudStorage()

	// Initialize Redis client
	initRedis()

	// Initialize NATS connection
	initNATS()

	// Initialize Elasticsearch client
	initElastic()
}

func initGoogleCloudStorage() {
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
}

func initRedis() {
	redisCtx := context.Background()
	var redisErr error
	for attempt := 1; attempt <= 3; attempt++ {
		RedisClient, _, redisErr = connectToRedis(redisCtx)
		if redisErr == nil {
			log.Println("Redis: OK")
			break
		}
		log.Printf("Failed to connect to Redis (attempt %d): %v\n", attempt, redisErr)
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	if redisErr != nil {
		log.Fatalf("Failed to connect to Redis after multiple attempts: %v", redisErr)
	}
}

func initNATS() {
	natsOpts := nats.Options{
		Servers: []string{os.Getenv("NATS_ADDRESS")},
	}

	var err error
	NatsClient, err = natsOpts.Connect()
	if err != nil {
		log.Fatalf("Error connecting to NATS server: %v", err)
	}

	log.Println("connected to NATS")
	// defer NatsClient.Close()
}

// func initElastic() {
// 	// Read the ELASTICSEARCH_URL from the environment variable
// 	elasticsearchURL := os.Getenv("ELASTICSEARCH_URL")
// 	if elasticsearchURL == "" {
// 		log.Fatal("ELASTICSEARCH_URL environment variable is not set")
// 	}

// 	// Load the CA certificate
// 	caCertPool, e := loadCA("/certs/elasticsearch.crt")
// 	if e != nil {
// 		log.Fatalf("Failed to load CA certificate: %s", e)
// 	}

// 	// Initialize the Elasticsearch client with the URL from the environment variable
// 	cfg := elasticsearch.Config{
// 		Addresses: []string{
// 			elasticsearchURL,
// 		},
// 		Transport: &http.Transport{
// 			TLSClientConfig: &tls.Config{
// 				RootCAs: caCertPool,
// 			},
// 		},
// 	}

// 	var err error
// 	ElasticClient, err = elasticsearch.NewClient(cfg)
// 	if err != nil {
// 		log.Fatalf("Error creating the elastic client: %s", err)
// 	}
// 	log.Println("Created Elastic client")
// }

func initElastic() {
	// Read the ELASTICSEARCH_URL from the environment variable
	elasticsearchURL := os.Getenv("ELASTICSEARCH_URL")
	elasticsearchUsername := os.Getenv("ELASTICSEARCH_USERNAME")
	elasticsearchPassword := os.Getenv("ELASTICSEARCH_PASSWORD")

	if elasticsearchURL == "" {
		log.Fatal("ELASTICSEARCH_URL environment variable is not set")
	}

	// Initialize the Elasticsearch client with the URL from the environment variable
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticsearchURL,
		},
		Username: elasticsearchUsername,
		Password: elasticsearchPassword,
	}
	var err error
	ElasticClient, err = elasticsearch.NewClient(cfg) // Corrected line
	if err != nil {
		log.Fatalf("Error creating the elastic client: %s", err)
	}
	log.Println("Created Elastic client with basic authentication")
}

// func loadCA(caCertPath string) (*x509.CertPool, error) {
// 	caCert, err := os.ReadFile(caCertPath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	caCertPool := x509.NewCertPool()
// 	if !caCertPool.AppendCertsFromPEM(caCert) {
// 		return nil, fmt.Errorf("failed to append CA certificate")
// 	}
// 	return caCertPool, nil
// }

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		// AllowOrigins:     "*",
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "x-jwt,x-spotify-id,x-spotify-token,Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Access-Control-Allow-Credentials",
		ExposeHeaders:    "Content-Length,Access-Control-Allow-Headers,Access-Control-Allow-Credentials",
		AllowCredentials: true, // this way, front-end can get the cookie and send it back
	}))
	app.Use(logger.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	auth := app.Group("/auth")
	music := app.Group("/music")
	places := app.Group("/places")
	playlists := app.Group("/playlists")
	revenues := app.Group("/revenues")
	users := app.Group("/users")
	communications := app.Group("/comms")
	search := app.Group("/search")

	authroute.AuthRoutes(auth, StorageClient, RedisClient, NatsClient)
	placesroute.PlaceRoutes(places, NatsClient)
	musicroute.MusicRoutes(music, StorageClient, RedisClient)
	playlistroute.PlaylistRoutes(playlists)
	revenueroute.RevenueRoutes(revenues)
	userroute.UserRoutes(users, RedisClient, NatsClient)
	communicationroutes.CommunicationRoutes(communications, RedisClient, NatsClient)
	searchroute.SearchRoutes(search, ElasticClient, NatsClient)

	// // below protected. Require an Authorization Bearer <token> to access.
	// app.Use(authroute.JWTSignerMiddleware(os.Getenv("JWT_SECRET")))
	// // Example of a protected route
	// app.Get("/protected", func(c *fiber.Ctx) error {
	// 	return c.SendString("Welcome to the protected area!")
	// })
	log.Print("APP @ 3000 : OK")
	log.Fatal(app.Listen(":3000"))

	// Configure HTTPS
	// certFile := "/certs/localhost.pem"
	// keyFile := "/certs/localhost-key.pem"
	// if err := app.ListenTLS(":3000", certFile, keyFile); err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }

	defer NatsClient.Close()

}

func connectToRedis(ctx context.Context) (*redis.Client, *redis.Options, error) {
	options := &redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: "",
		DB:       0,
	}
	client := redis.NewClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, options, err
	}
	return client, options, nil
}
