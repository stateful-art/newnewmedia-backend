package service

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type MailerService struct {
	smtpService *SMTPService
}

func NewMailerService(smtpService *SMTPService) *MailerService {
	return &MailerService{smtpService: smtpService}
}

func (ms *MailerService) StartVerification(c *fiber.Ctx, redisClient *redis.Client) error {
	email := c.Query("email")
	verificationLink, err := ms.smtpService.SendVerificationMail(email, redisClient)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"message": "lmao its not working",
		})
	}

	fmt.Printf("Here is your link: %v", verificationLink)

	return c.JSON(fiber.Map{
		"message": "Sent!",
	})
}

func (ms *MailerService) VerifyEmail(c *fiber.Ctx, redisClient *redis.Client) error {
	token := c.Query("token")
	email, err := redisClient.Get(context.Background(), token).Result()
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Error"})
	}
	fmt.Println(email)
	fmt.Println("Now you can update emailVerified column to true")

	_, err = redisClient.Del(context.Background(), token).Result()
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Error"})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Your email, %s is now verified", email),
	})
}
