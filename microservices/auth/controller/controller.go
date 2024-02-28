package controller

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	service "newnewmedia.com/microservices/auth/service" // Import your service package
)

type SpotifyToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	SpotifyID    string `json:"spotify_id"`
}

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
func SpotifyCallback(c *fiber.Ctx, redisClient *redis.ClusterClient) error {
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

	if spotifyToken == nil {
		// Handle nil token
		return errors.New("received nil Spotify token")
	}

	// Store Spotify token in Redis
	if err := storeSpotifyToken(redisClient, spotifyToken); err != nil {
		// Log error
		log.Println("Error storing Spotify token in Redis:", err)

		// Return error
		return err
	}

	// Redirect the client to the specified URL
	redirectURL := fmt.Sprintf("%s/?code=%s&refresh=%s&expire=%s", os.Getenv("TEST_WEBAPP_ORIGIN"), spotifyToken.AccessToken, spotifyToken.RefreshToken, strconv.FormatInt(spotifyToken.ExpiresIn, 10))
	c.Redirect(redirectURL, fiber.StatusSeeOther)

	// Return nil to indicate success
	return nil
}

func (st *SpotifyToken) MarshalBinary() ([]byte, error) {
	expiresInBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(expiresInBytes, uint64(st.ExpiresIn))

	// Concatenate all fields into a byte slice
	tokenBytes := []byte(st.AccessToken + st.RefreshToken + st.SpotifyID)

	// Combine expiresInBytes and tokenBytes
	return append(expiresInBytes, tokenBytes...), nil
}

func (st *SpotifyToken) UnmarshalBinary(data []byte) error {
	// Implement UnmarshalBinary if needed
	return nil
}

// storeSpotifyToken stores the Spotify token in Redis
func storeSpotifyToken(redisClient *redis.ClusterClient, spotifyToken *service.SpotifyToken) error {
	// Check if redisClient is nil
	if redisClient == nil {
		return errors.New("redisClient is nil")
	}

	// Convert Spotify token to JSON
	tokenJSON, err := json.Marshal(spotifyToken)
	if err != nil {
		// Handle error
		return err
	}

	// Set the Spotify token in Redis with SpotifyID as the key
	key := "spotify:" + spotifyToken.SpotifyID
	if err := redisClient.Set(context.Background(), key, tokenJSON, 0).Err(); err != nil {
		// Log error
		log.Println("Error setting Spotify token in Redis:", err)

		return err
	}

	return nil
}

// func storeSpotifyToken(redisClient *redis.Client, spotifyToken *service.SpotifyToken) error {
// 	// Convert Spotify token to JSON
// 	tokenJSON, err := json.Marshal(spotifyToken)
// 	if err != nil {
// 		return err
// 	}

// 	// Set the Spotify token in Redis with SpotifyID as the key
// 	key := "spotify:" + spotifyToken.SpotifyID
// 	if err := redisClient.Set(context.Background(), key, tokenJSON, 0).Err(); err != nil {
// 		return err
// 	}

// 	return nil
// }
