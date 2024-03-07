package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	service "newnew.media/microservices/communication/service"
)

type CommunicationController struct {
	mailerService *service.MailerService
	smtpService   *service.SMTPService
}

func NewCommunicationController(mailerService *service.MailerService, smtpService *service.SMTPService) *CommunicationController {
	return &CommunicationController{mailerService: mailerService, smtpService: smtpService}
}

func (cc *CommunicationController) StartVerification(c *fiber.Ctx, redisClient *redis.Client) error {
	email := c.Query("email")

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Email parameter is required"})
	}

	// Fetch the audio file path for the given song ID
	_, err := cc.smtpService.SendVerificationMail(email, redisClient)
	if err != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Could not send the email. Please try again later."})
	}
	return c.SendStatus(fiber.StatusAccepted)
}

func (cc *CommunicationController) VerifyEmail(c *fiber.Ctx, redisClient *redis.Client) error {
	error := cc.mailerService.VerifyEmail(c, redisClient)
	if error != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Could not send email"})
	}
	return nil
}
