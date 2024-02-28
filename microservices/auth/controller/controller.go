package controller

import (
	"fmt"
	"os"
	"strconv"

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
	return c.Redirect(authURL)

}

// // SpotifyCallback handles the callback from Spotify after user authorization
// func SpotifyCallback(c *fiber.Ctx) error {
// 	code := c.Query("code")
// 	if code == "" {
// 		// Handle error
// 		return fiber.ErrBadRequest
// 	}
// 	spotifyToken, err := service.HandleSpotifyCallback(c, code)
// 	if err != nil {
// 		// Handle error
// 		return err
// 	}
// 	return c.JSON(spotifyToken)
// }

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

	// Redirect the client to the specified URL

	// Set each token value as a separate item in the response header
	redirectURL := fmt.Sprintf("%s/?code=%s&refresh=%s&expire=%s", os.Getenv("TEST_WEBAPP_ORIGIN"), spotifyToken.AccessToken, spotifyToken.RefreshToken, strconv.FormatInt(spotifyToken.ExpiresIn, 10))
	c.Redirect(redirectURL, fiber.StatusSeeOther)
	// Return nil to indicate success
	return nil
}
