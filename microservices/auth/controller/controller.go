package controller

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	service "newnew.media/microservices/auth/service" // Import your service package
	userDTO "newnew.media/microservices/user/dto"
)

type AuthController struct {
	emailAuthService   *service.EmailAuthService
	spotifyAuthService *service.SpotifyAuthService
}

func NewAuthController(emailAuthService *service.EmailAuthService, spotifyAuthService *service.SpotifyAuthService) *AuthController {
	return &AuthController{emailAuthService: emailAuthService, spotifyAuthService: spotifyAuthService}
}

// SpotifyLogin initiates the Spotify OAuth 2.0 authentication flow
func (ac *AuthController) SpotifyLogin(c *fiber.Ctx) error {
	authURL, err := ac.spotifyAuthService.ConnectSpotify()
	if err != nil {
		// Handle error
		return err
	}
	return c.Redirect(authURL)

}

// SpotifyCallback handles the callback from Spotify after user authorization
func (ac *AuthController) SpotifyCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return fiber.ErrBadRequest
	}
	spotifyToken, err := ac.spotifyAuthService.HandleSpotifyCallback(c, code)
	if err != nil {
		return err
	}

	if spotifyToken == nil {
		return errors.New("received nil Spotify token")
	}

	// Redirect the client to the specified URL
	log.Println("redirectin here with Expiry of 0")
	redirectURL := fmt.Sprintf("%s/?code=%s&refresh=%s", os.Getenv("WEBAPP_ORIGIN"), spotifyToken.AccessToken, spotifyToken.RefreshToken)
	c.Redirect(redirectURL, fiber.StatusSeeOther)

	return nil
}

// EMAIL & Password register & logins

// EmailRegistration handles the email registration process.
func (ac *AuthController) EmailRegistration(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user := userDTO.CreateUserRequest{
		Email:     email,
		Password:  password,
		SpotifyID: "",
	}

	status, err := ac.emailAuthService.RegisterUser(user)
	if err != nil {
		return c.JSON(fiber.Map{"status": status, "error": err.Error()})
	}
	return nil
}

// EmailLogin handles the email login process.
func (ac *AuthController) EmailLogin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	token, err := ac.emailAuthService.LoginUser(email, password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	} else {

		return c.JSON(userDTO.LoginUserResponse{Email: email, Token: token, RefreshToken: ""})
	}

}
