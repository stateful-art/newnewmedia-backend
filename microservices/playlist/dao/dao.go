package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Song struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Artist    primitive.ObjectID `json:"artist,omitempty" bson:"artist_id,omitempty"`
	PlayCount int16              `json:"play_count,omitempty" bson:"play_count,omitempty"`
}

type PlaylistType string

const (
	Private PlaylistType = "private"
	Public  PlaylistType = "public"
)

type RevenueSharingModel string

const (
	CollectiveSharing RevenueSharingModel = "collective"
	IndividualSharing RevenueSharingModel = "individual"
)

type Playlist struct {
	ID                   primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	Name                 string              `json:"name,omitempty" bson:"name,omitempty"`
	Description          string              `json:"description,omitempty" bson:"description,omitempty"`
	Owner                primitive.ObjectID  `json:"owner,omitempty" bson:"owner,omitempty"`
	Type                 PlaylistType        `json:"type,omitempty" bson:"type,omitempty"`
	RevenueSharingModel  RevenueSharingModel `json:"revenue_sharing_model,omitempty" bson:"revenue_sharing_model,omitempty"`
	RevenueCutPercentage float64             `json:"revenue_cut_percentage,omitempty" bson:"revenue_cut_percentage,omitempty"`
	Songs                []Song              `json:"songs,omitempty" bson:"songs,omitempty"`
	CreatedAt            time.Time           `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt            time.Time           `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
