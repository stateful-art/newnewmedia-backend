package dao

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type Place struct {
	ID          string   `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string   `json:"name"`
	Location    Location `json:"location"`
	Description string   `json:"description"`
}

type RecentPlayItem struct {
	SongID     primitive.ObjectID `bson:"song_id"`
	SongName   string             `bson:"song_name"`
	ArtistName string             `bson:"artist_name"`
	PlaceName  string             `bson:"place_name"`
	Duration   int                `bson:"duration"`
	PlaceID    primitive.ObjectID `bson:"place_id"`
	// Other recent play item attributes...
}

type RecentPlays struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Place primitive.ObjectID `bson:"place_id"` // For place's recent plays
	Items []RecentPlayItem   `bson:"items"`
}
