package dao

import "go.mongodb.org/mongo-driver/bson/primitive"

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

type Role string

const (
	Audience Role = "audience"
	Artist   Role = "artist"
	Place    Role = "place"
	Admin    Role = "admin"
	Crew     Role = "crew"
)

type UserRole struct {
	UserID primitive.ObjectID `bson:"user_id"`
	Role   Role               `bson:"role"`
}

type FavoriteGenre struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserID  primitive.ObjectID `bson:"user_id"`
	GenreID primitive.ObjectID `bson:"genre_id"`
}

type FavoritePlace struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserID  primitive.ObjectID `bson:"user_id"`
	PlaceID primitive.ObjectID `bson:"place_id"`
}

type RecentPlayItem struct {
	SongID     primitive.ObjectID `bson:"song_id"`
	SongName   string             `bson:"song_name"`
	ArtistName string             `bson:"artist_name"`
	PlaceName  string             `bson:"place_name"`
	Duration   int                `bson:"duration"`
	PlaceID    primitive.ObjectID `bson:"place_id"`
}

type RecentPlays struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	User  primitive.ObjectID `bson:"user_id"` // For audience's recent plays
	Items []RecentPlayItem   `bson:"items"`
}
