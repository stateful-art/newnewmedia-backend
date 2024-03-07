package dto

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

type CreateUserRequest struct {
	Email    string
	Password string
}

type CreateUserResponse struct {
	Email     string
	SpotifyID string
	Token     string
}
type Role string

const (
	Audience Role = "audience"
	Artist   Role = "artist"
	Place    Role = "place"
	Admin    Role = "admin"
	Crew     Role = "crew"
)

type UserRoles struct {
	UserID primitive.ObjectID `bson:"user_id"`
	Roles  []Role             `bson:"role"`
}

type UpdateRoleRequest struct {
	UserID string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Role   Role   `json:"role,omitempty" bson:"role,omitempty"`
}
