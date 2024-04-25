package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"newnew.media/microservices/playlist/dao"
)

type Song struct {
	ID        string `json:"id,omitempty" bson:"id,omitempty"`
	Name      string `json:"name,omitempty" bson:"name,omitempty"`
	Artist    string `json:"artist,omitempty" bson:"artist_id,omitempty"`
	PlayCount int16  `json:"play_count,omitempty" bson:"play_count,omitempty"`
}

type PlaylistType string

const (
	Private PlaylistType = "private"
	Public  PlaylistType = "public"
)

type ContentSource string

const (
	Local   ContentSource = "local"
	Spotify ContentSource = "spotify"
	Youtube ContentSource = "youtube"
)

type RevenueSharingModel string

const (
	CollectiveSharing RevenueSharingModel = "collective"
	IndividualSharing RevenueSharingModel = "individual"
)

type CreatePlaylist struct {
	Name                 string              `json:"name,omitempty" bson:"name,omitempty"`
	Description          string              `json:"description,omitempty" bson:"description,omitempty"`
	Owner                string              `json:"owner,omitempty" bson:"owner,omitempty"`
	Type                 PlaylistType        `json:"type,omitempty" bson:"type,omitempty"`
	Source               ContentSource       `json:"content_source,omitempty" bson:"content_source,omitempty"`
	RevenueSharingModel  RevenueSharingModel `json:"revenue_sharing_model,omitempty" bson:"revenue_sharing_model,omitempty"`
	RevenueCutPercentage float64             `json:"revenue_cut_percentage,omitempty" bson:"revenue_cut_percentage,omitempty"`
	Songs                []Song              `json:"songs,omitempty" bson:"songs,omitempty"`
	Url                  string              `json:"url,omitempty" bson:"url,omitempty"`
	Image                string              `json:"image,omitempty" bson:"image,omitempty"`
}

type GetPlaylist struct {
	ID                   string              `json:"_id,omitempty" bson:"_id,omitempty"`
	Name                 string              `json:"name,omitempty" bson:"name,omitempty"`
	Description          string              `json:"description,omitempty" bson:"description,omitempty"`
	Owner                string              `json:"owner,omitempty" bson:"owner,omitempty"`
	Type                 PlaylistType        `json:"type,omitempty" bson:"type,omitempty"`
	Source               ContentSource       `json:"content_source,omitempty" bson:"content_source,omitempty"`
	RevenueSharingModel  RevenueSharingModel `json:"revenue_sharing_model,omitempty" bson:"revenue_sharing_model,omitempty"`
	RevenueCutPercentage float64             `json:"revenue_cut_percentage,omitempty" bson:"revenue_cut_percentage,omitempty"`
	Songs                []Song              `json:"songs,omitempty" bson:"songs,omitempty"`
	Url                  string              `json:"url,omitempty" bson:"url,omitempty"`
	Image                string              `json:"image,omitempty" bson:"image,omitempty"`
	CreatedAt            string              `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt            string              `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type RevenueCalculateRequest struct {
	Playlist     dao.Playlist `json:"playlist"`
	TotalRevenue float64      `json:"totalRevenue"`
}

type PlaylistSongsUpdateRequest struct {
	PlaylistID primitive.ObjectID   `json:"playlist_id"`
	SongIDs    []primitive.ObjectID `json:"song_ids"`
}
