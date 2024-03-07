package communicationroutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	controller "newnew.media/microservices/communication/controller"
	service "newnew.media/microservices/communication/service"
)

func CommunicationRoutes(app fiber.Router, redisClient *redis.Client) {
	// Create an instance of PlaylistService with the PlaylistRepository
	smtpService := service.NewSMTPService()
	mailerService := service.NewMailerService(smtpService)

	// Pass the PlaylistService instance to the PlaylistController
	cc := controller.NewCommunicationController(mailerService, smtpService)

	app.Get("/send-verification-email", func(c *fiber.Ctx) error {
		return cc.StartVerification(c, redisClient)
	})

	app.Get("/verify-email", func(c *fiber.Ctx) error {
		return cc.VerifyEmail(c, redisClient)
	})

}
