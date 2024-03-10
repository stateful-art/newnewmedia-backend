package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user with basic information and connected Spotify/YouTube Music accounts
type User struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email          string
	Password       string
	City           string
	FavoriteGenres []string
	FavoritePlaces []string
	SpotifyID      string
	YouTubeID      string
	CreatedAt      time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	EmailSent      bool      `json:"email_sent,omitempty" bson:"email_sent,omitempty"`
	EmailVerified  bool      `json:"email_verified,omitempty" bson:"email_verified,omitempty"`
}

type CreateUserRequest struct {
	Email     string
	Password  string
	City      string
	SpotifyID string
}

type CreateUserResponse struct {
	Status bool
}

type LoginUserResponse struct {
	Email        string
	Token        string
	RefreshToken string
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
