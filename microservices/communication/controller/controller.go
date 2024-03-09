package controller

import (
	"github.com/gofiber/fiber/v2"
	service "newnew.media/microservices/communication/service"
)

type CommunicationController struct {
	mailerService *service.MailerService
}

func NewCommunicationController(mailerService *service.MailerService) *CommunicationController {
	return &CommunicationController{mailerService: mailerService}
}

func (cc *CommunicationController) StartVerification(c *fiber.Ctx) error {
	email := c.Query("email")

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Email parameter is required"})
	}

	// Fetch the audio file path for the given song ID
	_, err := cc.mailerService.SendVerificationEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Could not send the email. Please try again later."})
	}
	return c.SendStatus(fiber.StatusAccepted)
}

func (cc *CommunicationController) VerifyEmail(c *fiber.Ctx) error {
	error := cc.mailerService.VerifyEmail(c)
	if error != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Could not send email"})
	}
	return nil
}
