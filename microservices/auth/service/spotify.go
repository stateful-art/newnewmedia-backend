package service

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	// Import your DAO package where User struct is defined
)

var (
	spotifyClientID     = os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	spotifyRedirectURI  = os.Getenv("SPOTIFY_REDIRECT_URI")
	spotifyAuthURL      = os.Getenv("SPOTIFY_AUTH_URL")
	spotifyTokenURL     = os.Getenv("SPOTIFY_TOKEN_URL")
)

type SpotifyToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func spotifyOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     spotifyClientID,
		ClientSecret: spotifyClientSecret,
		RedirectURL:  spotifyRedirectURI,
		Scopes:       []string{spotify.ScopePlaylistReadPrivate, spotify.ScopeUserReadRecentlyPlayed},
		Endpoint: oauth2.Endpoint{
			AuthURL:  spotifyAuthURL,
			TokenURL: spotifyTokenURL,
		},
	}
}

// ConnectSpotify initiates the Spotify OAuth 2.0 authentication flow
func ConnectSpotify() (string, error) {
	spotifyConfig := spotifyOAuthConfig()
	// Generate authorization URL
	authURL := spotifyConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return authURL, nil
}

// // HandleSpotifyCallback handles the callback from Spotify after user authorization
// func HandleSpotifyCallback(c *fiber.Ctx, code string) (string, error) {
// 	spotifyConfig := spotifyOAuthConfig()

// 	// Exchange authorization code for access token
// 	token, err := spotifyConfig.Exchange(context.Background(), code)
// 	log.Println("access_token in spotify callback")
// 	log.Println(token)
// 	c.Locals("access_token", token)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Create Spotify client with the obtained access token
// 	client := spotifyConfig.Client(context.Background(), token)
// 	spotifyClient := spotify.NewClient(client)

// 	// Retrieve user's Spotify ID
// 	user, err := spotifyClient.CurrentUser()
// 	if err != nil {
// 		return "", err
// 	}

// 	fmt.Print(user.ID)

// 	// Create a new user with the Spotify ID
// 	newUser := dao.User{
// 		SpotifyID: user.ID,
// 		// Set other fields as needed
// 	}

// 	fmt.Print(newUser)
// 	// Save the user to the database or perform any other necessary actions

// 	return user.ID, nil
// }

func HandleSpotifyCallback(c *fiber.Ctx, code string) (*SpotifyToken, error) {
	spotifyConfig := spotifyOAuthConfig()

	// Exchange authorization code for access token
	token, err := spotifyConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	// Log the token
	log.Println("access_token in spotify callback")
	log.Println(token)

	// Calculate the expiration time in seconds
	expiresIn := int64(time.Until(token.Expiry).Seconds())

	// Create a SpotifyToken struct
	spotifyToken := &SpotifyToken{
		AccessToken:  token.AccessToken,
		ExpiresIn:    expiresIn,
		RefreshToken: token.RefreshToken,
	}
	// Convert the SpotifyToken struct to JSON
	spotifyTokenJSON, err := json.Marshal(spotifyToken)
	if err != nil {
		return nil, err
	}

	// Log the JSON representation of the SpotifyToken
	log.Println("SpotifyToken JSON:", string(spotifyTokenJSON))

	return spotifyToken, nil
}

// FetchRecentlyPlayedSongs fetches the user's recently played songs.
func FetchRecentlyPlayedSongs(token *oauth2.Token) ([]spotify.RecentlyPlayedItem, error) {
	// Create Spotify client with the obtained access token
	client := spotify.NewAuthenticator("").NewClient(token)

	// Retrieve user's recently played tracks
	tracks, err := client.PlayerRecentlyPlayed()
	if err != nil {
		return nil, err
	}

	return tracks, nil
}
