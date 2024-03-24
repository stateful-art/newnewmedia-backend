package userroutes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	controller "newnew.media/microservices/user/controller"
	dao "newnew.media/microservices/user/dao"
	repo "newnew.media/microservices/user/repository"
	service "newnew.media/microservices/user/service"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func UserRoutes(app fiber.Router, redisClient *redis.Client, natsClient *nats.Conn) {
	// Create an instance of PlaylistRepository
	userRepo := repo.NewUserRepository()
	// Create an instance of PlaylistService with the PlaylistRepository
	userService := service.NewUserService(userRepo, redisClient, natsClient)

	// Pass the PlaylistService instance to the PlaylistController
	userController := controller.NewUserController(userService)

	// Define routes using the controller methods
	app.Get("/", userController.GetUsers)
	app.Get("/:id", userController.GetUserByID)
	app.Get("/:spotify_id", userController.GetUserBySpotifyID)
	app.Get("/:youtube_id", userController.GetUserByYouTubeID)

	app.Post("/", userController.CreateUser)
	app.Patch("/add-role", CheckAdminRole, userController.AddRole)
	app.Patch("/remove-role", CheckAdminRole, userController.RemoveRole)

}

func CheckAdminRole(c *fiber.Ctx) error {
	// Extract the JWT token from the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Missing Authorization header"})
	}

	// Parse the JWT token
	token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
	}

	// Check if the token is valid and contains the necessary claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract the roles from the claims
		roles, ok := claims["roles"].([]interface{})
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "No roles found in token"})
		}

		// Check if the user has the Admin role
		for _, role := range roles {
			if role == string(dao.Admin) {
				return c.Next() // Proceed to the next handler if the user has the Admin role
			}
		}
	}

	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Insufficient permissions"})
}
