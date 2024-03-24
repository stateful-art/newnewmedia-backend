package dao

import "go.mongodb.org/mongo-driver/bson/primitive"

// PlaceDto is a struct that represents the data transfer object for the place microservice
type Location struct {
	Address   string  `json:"address"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Link struct {
	Platform string `json:"platform"`
	Url      string `json:"url"`
}

type Place struct {
	ID          string   `json:"id,omitempty" bson:"_id,omitempty"`
	Owner       string   `json:"owner,omitempty" bson:"owner,omitempty"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	SpotifyID   string   `json:"spotify_id"`
	Name        string   `json:"name"`
	Location    Location `json:"location"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Links       []Link   `json:"links"`
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
