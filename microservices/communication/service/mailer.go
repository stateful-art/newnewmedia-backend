package service

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

const REDIS_UNVERIFIED_EMAIL_PREFIX = "unverified"

type MailerService struct {
	smtpService *SMTPService
	natsClient  *nats.Conn
}

func NewMailerService(smtpService *SMTPService, natsClient *nats.Conn) *MailerService {
	return &MailerService{smtpService: smtpService, natsClient: natsClient}
}

func (ms *MailerService) StartVerification(c *fiber.Ctx, redisClient *redis.Client) error {
	email := c.Query("email")
	verificationLink, err := ms.smtpService.SendVerificationMail(email)
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
	email, err := redisClient.Get(context.Background(), fmt.Sprintf("%s:%s", REDIS_UNVERIFIED_EMAIL_PREFIX, token)).Result()
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Error"})
	}
	fmt.Println(email)
	fmt.Println("Now you can update emailVerified column to true")

	_, err = redisClient.Del(context.Background(), fmt.Sprintf("%s:%s", REDIS_UNVERIFIED_EMAIL_PREFIX, token)).Result()
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Error"})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Your email, %s is now verified", email),
	})
}

func (ms *MailerService) SubscribeToUserRegisteredSubject() error {
	// Subscribe to user-registered subject
	_, err := ms.natsClient.Subscribe("user-registered", func(msg *nats.Msg) {
		email := string(msg.Data)
		// Send verification email using SMTP service
		_, err := ms.smtpService.SendVerificationMail(email) // Pass nil for redisClient, as it's not needed here
		if err != nil {
			log.Printf("Failed to send verification email to %s: %v\n", email, err)
			return
		}
		// log.Printf("Verification email sent to %s. Verification link: %s\n", email, verificationLink)
		select {}
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to user-registered subject: %v", err)
	}
	// defer sub.Unsubscribe()
	return nil
}
