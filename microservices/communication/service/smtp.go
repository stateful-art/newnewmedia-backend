package service

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SMTPService struct {
}

func NewSMTPService() *SMTPService {
	return &SMTPService{}
}

func (smtps *SMTPService) SendVerificationMail(email string, redisClient *redis.Client) (string, error) {

	backendOrigin := os.Getenv("BACKEND_ORIGIN")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_MAIL")
	password := os.Getenv("SMTP_PASSWORD")

	log.Println(host)
	log.Println(port)
	log.Println(from)
	log.Println(password)

	subject := "Please verify your email"
	auth := smtp.PlainAuth("", from, password, host)

	verificationKey := fmt.Sprintf(
		"%x", sha256.Sum256([]byte(email + "-" + uuid.New().String())[:]),
	)

	verificationLinkBase := fmt.Sprintf("%s/comms/verify-email?token=", backendOrigin)

	// body := fmt.Sprintf(`
	// <html>
	// <a href="%v%v" target="_blank">CLICK</a>
	// <p>This link will expire in 24 hours.</p>
	// </html>
	// `, verificationLinkBase, verificationKey)

	body := fmt.Sprintf(`
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
				<p>Please do not reply to this email.</p>
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

	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}

	c, err := smtp.Dial(host + ":" + port)
	if err != nil {
		return "", err
	}
	c.StartTLS(tlsconfig)

	if err = c.Auth(auth); err != nil {
		fmt.Println(err)
		return "", err
	}

	if err = c.Mail(from); err != nil {
		return "", err
	}

	if err = c.Rcpt(email); err != nil {
		return "", err
	}

	w, err := c.Data()
	if err != nil {
		return "", err
	}

	_, err = w.Write([]byte(
		fmt.Sprintf("MIME-Version: %v\r\n", "1.0") +
			fmt.Sprintf("Content-type: %v\r\n", "text/html; charset=UTF-8") +
			fmt.Sprintf("From: %v\r\n", from) +
			fmt.Sprintf("To: %v\r\n", email) +
			fmt.Sprintf("Subject: %v\r\n", subject) +
			fmt.Sprintf("%v\r\n", body),
	))

	if err != nil {
		return "", err
	}

	statusCMD := redisClient.Set(context.Background(),
		verificationKey,
		fmt.Sprintf("%v", email),
		time.Duration(time.Hour*24),
	)
	if statusCMD.Err() != nil {
		return "", statusCMD.Err()
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	c.Quit()

	return fmt.Sprintf("%v%v", verificationLinkBase, verificationKey), nil

}
