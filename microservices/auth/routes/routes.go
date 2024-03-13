package authroutes

import (
	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	controller "newnew.media/microservices/auth/controller"
	as "newnew.media/microservices/auth/service" // Import your service package
	ur "newnew.media/microservices/user/repository"
	us "newnew.media/microservices/user/service"
)

func AuthRoutes(app fiber.Router, storageClient *storage.Client, redisClient *redis.Client, natsClient *nats.Conn) {
	userRepository := &ur.UserRepository{}
	userService := us.NewUserService(userRepository, redisClient, natsClient)

	emailAuthService := as.NewEmailAuthService(userService)

	spotifyAuthService := as.NewSpotifyAuthService(natsClient, redisClient, userService, nil)
	spotifyAuthService.StartTokenRefresher()

	authController := controller.NewAuthController(emailAuthService, spotifyAuthService)

	app.Get("/spotify", authController.SpotifyLogin) // Initiate Spotify login
	app.Get("/spotify/callback", authController.SpotifyCallback)
	app.Post("/email/register", authController.EmailRegistration)
	app.Post("/email/login", authController.EmailLogin)
}

// NewAuthMiddleware creates a new JWT middleware with the given secret key.
func JWTSignerMiddleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(secret),
	})
}
