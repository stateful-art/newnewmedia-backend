package communicationroutes

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	controller "newnew.media/microservices/communication/controller"
	service "newnew.media/microservices/communication/service"
	userService "newnew.media/microservices/user/service"
)

func CommunicationRoutes(app fiber.Router, redisClient *redis.Client, natsClient *nats.Conn) {
	var domain string = os.Getenv("EMAIL_SENDER_DOMAIN")
	var key string = os.Getenv("MAILGUN_APIKEY")

	mailgun := mailgun.NewMailgun(domain, key)
	mailerService := service.NewMailerService(mailgun, natsClient, redisClient, &userService.UserService{})

	// Subscribe to user-registered subject
	if err := mailerService.SubscribeToUserRegisteredSubject(); err != nil {
		log.Fatalf("Failed to subscribe to user-registered subject: %v", err)
	}

	// Pass the PlaylistService instance to the PlaylistController
	cc := controller.NewCommunicationController(mailerService)

	app.Get("/send-verification-email", cc.StartVerification)
	app.Get("/verify-email", cc.VerifyEmail)

}
