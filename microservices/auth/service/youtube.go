package service

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	youtubeClientID     = os.Getenv("YOUTUBE_CLIENT_ID")
	youtubeClientSecret = os.Getenv("YOUTUBE_CLIENT_SECRET")
	youtubeRedirectURI  = os.Getenv("YOUTUBE_REDIRECT_URI")
	youtubeAuthURL      = os.Getenv("YOUTUBE_AUTH_URI")
	youtubeTokenURL     = os.Getenv("YOUTUBE_TOKEN_URI")
)

// ConnectYouTubeMusic initiates the YouTube Music OAuth 2.0 authentication flow
func ConnectYouTubeMusic() (string, error) {
	// Set up OAuth2 config for YouTube Music
	youtubeOAuthConfig := &oauth2.Config{
		ClientID:     youtubeClientID,
		ClientSecret: youtubeClientSecret,
		RedirectURL:  youtubeRedirectURI,
		Scopes:       []string{youtube.YoutubeReadonlyScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  youtubeAuthURL,
			TokenURL: youtubeTokenURL,
		},
	}

	// Generate authorization URL
	authURL := youtubeOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)

	return authURL, nil
}

// HandleYouTubeMusicCallback handles the callback from YouTube Music after user authorization
func HandleYouTubeMusicCallback(code string) error {
	// Set up OAuth2 config for YouTube Music
	youtubeOAuthConfig := &oauth2.Config{
		ClientID:     youtubeClientID,
		ClientSecret: youtubeClientSecret,
		RedirectURL:  youtubeRedirectURI,
		Scopes:       []string{youtube.YoutubeReadonlyScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  youtubeAuthURL,
			TokenURL: youtubeTokenURL,
		},
	}

	// Exchange authorization code for access token
	token, err := youtubeOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	youtubeService, err := youtube.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	// Example: Retrieve user's playlists
	playlists, err := youtubeService.Playlists.List([]string{"snippet"}).Do()
	if err != nil {
		return err
	}

	// Print the user's playlists
	fmt.Println("Playlists:")
	for _, playlist := range playlists.Items {
		fmt.Printf("  - %s\n", playlist.Snippet.Title)
	}

	return nil
}
