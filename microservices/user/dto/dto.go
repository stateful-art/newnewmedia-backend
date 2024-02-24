package dto

// User represents a user with basic information and connected Spotify/YouTube Music accounts
type User struct {
	Email          string
	Password       string
	City           string
	FavoriteGenres []string
	FavoritePlaces []string
	SpotifyID      string
	YouTubeID      string
}
