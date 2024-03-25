package dao

import "go.mongodb.org/mongo-driver/bson/primitive"

type GeometryType string

const (
	Point              GeometryType = "Point"
	LineString         GeometryType = "LineString"
	Polygon            GeometryType = "Polygon"
	MultiPoint         GeometryType = "MultiPoint"
	MultiLineString    GeometryType = "MultiLineString"
	MultiPolygon       GeometryType = "MultiPolygon"
	GeometryCollection GeometryType = "GeometryCollection"
)

// Location struct with a GeometryType field
type Location struct {
	Type        GeometryType `json:"type"`
	Coordinates []float64    `json:"coordinates"`
}
type Link struct {
	Platform string `json:"platform"`
	Url      string `json:"url"`
}

type Place struct {
	ID        string   `json:"id,omitempty" bson:"_id,omitempty"`
	Owner     string   `json:"owner,omitempty" bson:"owner,omitempty"`
	Email     string   `json:"email"`
	Phone     string   `json:"phone"`
	SpotifyID string   `json:"spotify_id" bson:"spotify_id"`
	Name      string   `json:"name"`
	Location  Location `json:"location"`
	City      string   `json:"city"`
	Country   string   `json:"country"`

	Description string `json:"description"`
	Image       string `json:"image"`
	Links       []Link `json:"links"`
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
	Place primitive.ObjectID `bson:"place_id"` // For place's recent plays
	Items []RecentPlayItem   `bson:"items"`
}
