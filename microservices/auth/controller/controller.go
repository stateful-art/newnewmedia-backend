package controller

import (
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
	c.Set("X-Spotify-Access-Token", spotifyToken.AccessToken)
	c.Set("X-Spotify-Expires-In", strconv.FormatInt(spotifyToken.ExpiresIn, 10))
	c.Set("X-Spotify-Refresh-Token", spotifyToken.RefreshToken)

	c.Redirect(os.Getenv("TEST_WEBAPP_ORIGIN"), fiber.StatusSeeOther)
	// Return nil to indicate success
	return nil
}
