package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"encoding/base64"

	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	userDTO "newnew.media/microservices/user/dto"
	userService "newnew.media/microservices/user/service"
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
	AccessToken  string    `json:"access_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	RefreshToken string    `json:"refresh_token"`
	SpotifyID    string    `json:"spotify_id"`
	Email        string    `json:"spotify_email"`
}

var scopes = []string{spotify.ScopePlaylistReadCollaborative,
	spotify.ScopePlaylistReadPrivate,
	spotify.ScopeUserReadCurrentlyPlaying,
	spotify.ScopeUserReadRecentlyPlayed,
	spotify.ScopeUserFollowRead,
	spotify.ScopeUserReadEmail,
	spotify.ScopeUserReadPlaybackState,
}

func spotifyOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     spotifyClientID,
		ClientSecret: spotifyClientSecret,
		RedirectURL:  spotifyRedirectURI,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  spotifyAuthURL,
			TokenURL: spotifyTokenURL,
		},
	}
}

type SpotifyAuthService struct {
	natsClient  *nats.Conn
	redisClient *redis.Client
	userService *userService.UserService
	config      *oauth2.Config
	mu          sync.Mutex               // Mutex for synchronizing access to token
	tokens      map[string]*oauth2.Token // Map to store tokens for multiple users
}

func NewSpotifyAuthService(natsClient *nats.Conn, redisClient *redis.Client, userService *userService.UserService, config *oauth2.Config) *SpotifyAuthService {
	return &SpotifyAuthService{natsClient: natsClient, redisClient: redisClient, userService: userService, config: spotifyOAuthConfig(), tokens: make(map[string]*oauth2.Token)}
}

// ConnectSpotify initiates the Spotify OAuth 2.0 authentication flow
func (sas *SpotifyAuthService) ConnectSpotify() (string, error) {
	// spotifyConfig := sas.config
	// Generate authorization URL
	authURL := sas.config.AuthCodeURL("state", oauth2.AccessTypeOffline)

	return authURL, nil
}
func (sas *SpotifyAuthService) HandleSpotifyCallback(c *fiber.Ctx, code string) (*SpotifyToken, error) {
	// Exchange authorization code for access token
	token, err := sas.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	user, _ := sas.getCurrentUser(token)
	if err != nil {
		// Handle error
		log.Println("error parsing expiry on received token from spotify..")
		panic(err)
	}

	// Create a SpotifyToken struct
	spotifyToken := &SpotifyToken{
		AccessToken:  token.AccessToken,
		ExpiresAt:    token.Expiry,
		RefreshToken: token.RefreshToken,
		SpotifyID:    user.ID,
	}

	// Store token for future refresh
	sas.mu.Lock()
	defer sas.mu.Unlock()
	sas.tokens[user.ID] = token

	// Store Spotify token in Redis
	if err := sas.storeSpotifyToken(spotifyToken); err != nil {
		log.Println("Error storing Spotify token in Redis:", err)
		return nil, err
	}

	newUser := userDTO.CreateUserRequest{
		SpotifyID: spotifyToken.SpotifyID,
		Email:     "",
		Password:  "",
		City:      "",
	}

	// if already a spotify logged-in user, just log them in without creating a user for them.
	if err := sas.userService.CheckUserExists(newUser); err != nil {
		log.Printf("This spotify user [ %s ] is already on our platform", newUser.SpotifyID)
		return spotifyToken, nil
	}

	log.Printf("creating user registered with ")
	error := sas.userService.CreateUser(newUser)

	if error != nil {
		return nil, err
	} else {
		log.Printf("new user created with spotify login...")

		return spotifyToken, nil
	}

}

func (sas *SpotifyAuthService) GetToken(spotifyID string) (*oauth2.Token, error) {
	sas.mu.Lock()
	defer sas.mu.Unlock()
	token, ok := sas.tokens[spotifyID]
	if !ok {
		return nil, fmt.Errorf("token not found for Spotify ID: %s", spotifyID)
	}
	return token, nil
}

func (sas *SpotifyAuthService) refreshToken() {
	ticker := time.NewTicker(time.Minute) // Check every minute
	defer ticker.Stop()

	for range ticker.C {
		log.Printf("Just checking tokens.. now >> %s", time.Now().String())

		sas.mu.Lock()
		for spotifyID := range sas.tokens {
			token, err := sas.getRefreshToken(spotifyID)
			if err != nil {
				log.Printf("Error getting refresh token for Spotify ID %s: %v", spotifyID, err)
				continue
			}

			expiry, err := sas.getTokenExpiry(spotifyID)
			if err != nil {
				log.Printf("Error getting token expiry for Spotify ID %s: %v", spotifyID, err)
				continue
			}

			// Only refresh token if it's about to expire
			if token == nil || time.Until(expiry) < time.Minute {
				newAccessToken, newRefreshToken, _ := sas.refreshAccessToken(token.RefreshToken)

				var spotifyToken = SpotifyToken{
					AccessToken:  newAccessToken,
					RefreshToken: newRefreshToken,
					ExpiresAt:    time.Now().Add(time.Hour),
					SpotifyID:    spotifyID,
					Email:        "", // Add email if available
				}

				data, err := json.Marshal(spotifyToken)
				if err != nil {
					log.Println("hey")
				}

				var token oauth2.Token
				if err := json.Unmarshal(data, &token); err != nil {
					log.Print("hello")
				}

				sas.tokens[spotifyID] = &token

				// Update token in Redis
				if err := sas.storeSpotifyToken(&spotifyToken); err != nil {
					log.Printf("Error updating token in Redis for Spotify ID %s: %v", spotifyID, err)
					continue
				}
			}
		}
		sas.mu.Unlock()
	}
}

func (sas *SpotifyAuthService) getRefreshToken(spotifyID string) (*oauth2.Token, error) {
	key := "spotify:" + spotifyID
	tokenJSON, err := sas.redisClient.Get(context.Background(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("a spotify token not found for Spotify ID: %s", spotifyID)
		}
		return nil, err
	}

	var token oauth2.Token
	if err := json.Unmarshal(tokenJSON, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (sas *SpotifyAuthService) getTokenExpiry(spotifyID string) (time.Time, error) {
	key := "spotify:" + spotifyID
	tokenJSON, err := sas.redisClient.Get(context.Background(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return time.Time{}, fmt.Errorf("a spotify token not found for Spotify ID: %s", spotifyID)
		}
		return time.Time{}, err
	}

	var token struct {
		AccessToken  string    `json:"access_token"`
		ExpiresAt    time.Time `json:"expires_at"`
		RefreshToken string    `json:"refresh_token"`
		SpotifyID    string    `json:"spotify_id"`
		Email        string    `json:"spotify_email"`
	}
	if err := json.Unmarshal(tokenJSON, &token); err != nil {
		return time.Time{}, err
	}
	return token.ExpiresAt, nil
}

func (sas *SpotifyAuthService) StartTokenRefresher() {
	go sas.refreshToken()
}

// storeSpotifyToken stores the Spotify token in Redis
func (sas *SpotifyAuthService) storeSpotifyToken(spotifyToken *SpotifyToken) error {
	// Check if redisClient is nil
	if sas.redisClient == nil {
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
	log.Printf("added %s key to redis.", key)

	if err := sas.redisClient.Set(context.Background(), key, tokenJSON, 0).Err(); err != nil {
		// Log error
		log.Println("Error setting Spotify token in Redis:", err)

		return err
	}

	return nil
}

func (sas *SpotifyAuthService) getCurrentUser(token *oauth2.Token) (*spotify.PrivateUser, error) {
	client := sas.config.Client(context.Background(), token)
	spotifyClient := spotify.NewClient(client)

	// Retrieve user's Spotify ID
	user, err := spotifyClient.CurrentUser()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (sas *SpotifyAuthService) refreshAccessToken(refreshToken string) (string, string, error) {
	data := "grant_type=refresh_token&refresh_token=" + refreshToken

	auth := spotifyClientID + ":" + spotifyClientSecret
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(spotifyTokenURL)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/x-www-form-urlencoded")
	req.Header.Set("Authorization", basicAuth)
	req.SetBodyString(data)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		log.Println("Error making POST request:", err)
		return "", "", err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		log.Println("Error: Request failed with status code:", resp.StatusCode())
		return "", "", fmt.Errorf("request failed with status code: %d", resp.StatusCode())
	}

	var response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.Println("Error decoding response body:", err)
		return "", "", err
	}

	return response.AccessToken, response.RefreshToken, nil
}
