package controller

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	service "newnewmedia.com/microservices/auth/service" // Import your service package
)

// SpotifyLogin initiates the Spotify OAuth 2.0 authentication flow
func SpotifyLogin(c *fiber.Ctx) error {
	authURL, err := service.ConnectSpotify()
	if err != nil {
		// Handle error
		return err
	}
	log.Print(authURL)
	return c.Redirect(authURL)
}

// SpotifyCallback handles the callback from Spotify after user authorization
func SpotifyCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		// Handle error
		return fiber.ErrBadRequest
	}
	spotifyToken, err := service.HandleSpotifyCallback(c, code)
	if err != nil {
		// Handle error
		return err
	}
	response := fmt.Sprintf("Spotify authentication successful! User's access token: %s", spotifyToken.AccessToken)
	return c.SendString(response)
}
