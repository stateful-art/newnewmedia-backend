package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"

	"github.com/google/uuid"
	userDTO "newnew.media/microservices/user/dto"
	userService "newnew.media/microservices/user/service"
)

const REDIS_UNVERIFIED_EMAIL_PREFIX = "unverified"

type MailerService struct {
	mailgun     mailgun.Mailgun
	natsClient  *nats.Conn
	redisClient *redis.Client
	userService *userService.UserService
}

func NewMailerService(mailgun mailgun.Mailgun, natsClient *nats.Conn, redisClient *redis.Client, userService *userService.UserService) *MailerService {
	return &MailerService{mailgun: mailgun, natsClient: natsClient, redisClient: redisClient, userService: userService}
}

func (ms *MailerService) StartVerification(c *fiber.Ctx, redisClient *redis.Client) error {
	email := c.Query("email")
	verificationLink, err := ms.SendVerificationEmail(email)
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

func (ms *MailerService) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	email, err := ms.redisClient.Get(context.Background(), fmt.Sprintf("%s:%s", REDIS_UNVERIFIED_EMAIL_PREFIX, token)).Result()
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Already verified."})
	}

	log.Printf("fetching user by their email >> [ %s ]\n", email)
	user, err := ms.userService.GetUserByEmail(email)
	if err != nil {
		log.Printf("Error getting user by email [ %s ] \n", email)
	}
	user.EmailVerified = true

	erro := ms.userService.UpdateUser(user.ID, user)
	if erro != nil {
		log.Printf("Error updating EmailVerified field for user >>  [ %s ] \n", user.ID.Hex())
	}

	rolerror := ms.userService.AddRole(user.ID, userDTO.Audience)
	if rolerror != nil {
		return c.Status(401).JSON(fiber.Map{"error": rolerror.Error()})
	}

	_, err = ms.redisClient.Del(context.Background(), fmt.Sprintf("%s:%s", REDIS_UNVERIFIED_EMAIL_PREFIX, token)).Result()
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
		// Send verification email using mailgun
		_, err := ms.SendVerificationEmail(email) // Pass nil for redisClient, as it's not needed here
		if err != nil {
			log.Printf("Failed to send verification email to %s: %v\n", email, err)
			return
		}
		select {}
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to user-registered subject: %v", err)
	}
	// defer sub.Unsubscribe()
	return nil
}

func (ms *MailerService) SendVerificationEmail(recipient string) (string, error) {

	// Create an instance of the Mailgun Client

	//When you have an EU-domain, you must specify the endpoint:
	ms.mailgun.SetAPIBase("https://api.eu.mailgun.net/v3")

	sender := os.Getenv("EMAIL_SENDER_ADDRESS")
	subject := "Please verify your email address"
	BACKEND_ORIGIN := os.Getenv("BACKEND_ORIGIN")

	// generate a token
	verificationKey := fmt.Sprintf(
		"%x", sha256.Sum256([]byte(recipient + "-" + uuid.New().String())[:]),
	)

	verificationLinkBase := fmt.Sprintf("%s/comms/verify-email?token=", BACKEND_ORIGIN)

	body := ms.generateBody(verificationLinkBase, verificationKey)
	message := ms.mailgun.NewMessage(sender, subject, "", recipient)
	message.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, id, err := ms.mailgun.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	// Write to REDIS.
	statusCMD := ms.redisClient.Set(context.Background(),
		fmt.Sprintf("%s:%s", REDIS_UNVERIFIED_EMAIL_PREFIX, verificationKey),
		fmt.Sprintf("%v", recipient),
		time.Duration(time.Hour*24),
	)

	if statusCMD.Err() != nil {
		return "", statusCMD.Err()
	}

	return fmt.Sprintf("%v%v", verificationLinkBase, verificationKey), nil

}

func (ms *MailerService) generateBody(verificationLinkBase string, verificationKey string) string {
	return fmt.Sprintf(`
	<html>
	<head>
		<style>
			.container {
				max-width: 600px;
				margin: 0 auto;
				padding: 20px;
				font-family: Arial, sans-serif;
			}
			.logo {
				display: block;
				margin: 0 auto;
			}
			.verification-btn {
				display: inline-block;
				background-color: #000;
				color: #fff; /* Set text color to white */
				padding: 10px 20px;
				text-decoration: none;
				border-radius: 5px;
				margin-top: 20px;
			}
			.verification-btn:hover {
				background-color: #222;
			}
			.verification-btn:hover,
			.verification-btn:active {
				color: #fff; /* Set text color to white on hover and active states */
			}
			.footer {
				margin-top: 30px;
				border-top: 1px solid #ccc;
				padding-top: 20px;
				font-size: 12px;
				color: #666;
			}
			.footer a {
				color: #666;
				text-decoration: none;
			}
			.start-logo {
				color: red;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<img class="logo" src="https://www.stateful.art/newnewmedia.png" alt="Logo" width="100">
			<h2>Email Verification</h2>
			<p>Thank you for signing up to newnew.media! <br> To verify your email address, please click the button below:</p>
			<a class="verification-btn" href="%v%v" target="_blank">Verify Email</a>
			<div class="footer">
				<p>
					newnew.media is developed by <a href="https://stateful.art">st<span class="start-logo">art</span> </a>
				  
					<br><br>
					<a href="mailto:contact@stateful.art">contact@stateful.art</a>
	
				</p>
			</div>
		</div>
	</body>
	</html>
`, verificationLinkBase, verificationKey)
}
