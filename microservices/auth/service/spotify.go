package service

import (
	"context"
	"fmt"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var (
	spotifyClientID     = os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	spotifyRedirectURI  = os.Getenv("SPOTIFY_REDIRECT_URI")
)

// ConnectSpotify initiates the Spotify OAuth 2.0 authentication flow
func ConnectSpotify() (string, error) {
	// Set up OAuth2 config for Spotify
	spotifyOAuthConfig := &oauth2.Config{
		ClientID:     spotifyClientID,
		ClientSecret: spotifyClientSecret,
		RedirectURL:  spotifyRedirectURI,
		Scopes:       []string{spotify.ScopePlaylistReadPrivate, spotify.ScopeUserReadRecentlyPlayed},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}

	// Generate authorization URL
	authURL := spotifyOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)

	return authURL, nil
}

// HandleSpotifyCallback handles the callback from Spotify after user authorization
func HandleSpotifyCallback(code string) error {
	// Set up OAuth2 config for Spotify
	spotifyOAuthConfig := &oauth2.Config{
		ClientID:     spotifyClientID,
		ClientSecret: spotifyClientSecret,
		RedirectURL:  spotifyRedirectURI,
		Scopes:       []string{spotify.ScopePlaylistReadPrivate, spotify.ScopeUserReadRecentlyPlayed},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}

	// Exchange authorization code for access token
	token, err := spotifyOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	// Create Spotify client with the obtained access token
	client := spotifyOAuthConfig.Client(context.Background(), token)
	spotifyClient := spotify.NewClient(client)

	// Retrieve user's recently played tracks as an example
	tracks, err := spotifyClient.PlayerRecentlyPlayed()
	if err != nil {
		return err
	}

	// Print the user's recently played tracks
	fmt.Println("Recently Played Tracks:")
	for _, track := range tracks {
		fmt.Printf("  - %s by %s\n", track.Track.Name, track.Track.Artists[0].Name)
	}

	return nil
}
