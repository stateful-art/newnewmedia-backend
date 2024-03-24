package controller

import (
	"errors"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	service "newnew.media/microservices/auth/service" // Import your service package
	utils "newnew.media/microservices/auth/utils"
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
	redirectURL := fmt.Sprintf("%s/?code=%s&refresh=%s", os.Getenv("WEBAPP_ORIGIN"), spotifyToken.AccessToken, spotifyToken.RefreshToken)
	c.Redirect(redirectURL, fiber.StatusSeeOther)

	return nil
}

// EMAIL & Password register & logins
func (ac *AuthController) EmailRegistration(c *fiber.Ctx) error {
	email := utils.TrimInput(c.FormValue("email"))
	password := utils.TrimInput(c.FormValue("password"))

	validEmail := utils.IsValidEmail(email)
	isValidPassword := utils.IsValidPassword(password)

	if !validEmail {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Not a valid email address."})
	}

	if !isValidPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Passwords should be at least 8-digits long and include lower-case, upper-case letters, special characters and a number."})
	}

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

func (ac *AuthController) EmailLogin(c *fiber.Ctx) error {
	email := utils.TrimInput(c.FormValue("email"))
	password := utils.TrimInput(c.FormValue("password"))

	validEmail := utils.IsValidEmail(email)

	if !validEmail {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Not a valid email address."})
	}

	token, err := ac.emailAuthService.LoginUser(email, password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	} else {

		return c.JSON(userDTO.LoginUserResponse{Email: email, Token: token, RefreshToken: ""})
	}

}
