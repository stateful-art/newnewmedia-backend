package dao

import "go.mongodb.org/mongo-driver/bson/primitive"

// TODO: Update this as "Song"
type Music struct {
	ID     string               `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string               `json:"name"`
	Artist string               `json:"artist"`
	Path   string               `json:"path"`
	Genres []primitive.ObjectID `bson:"genres"`
}

type Genre struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

type QueueSong struct {
	SongID       primitive.ObjectID `bson:"song_id"`
	SongName     string             `bson:"song_name"`
	ArtistName   string             `bson:"artist"`
	Duration     int                `bson:"duration"`
	PlaceID      primitive.ObjectID `bson:"place_id"`
	PlaylistName string             `bson:"playlist_name"`
}

type Queue struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserID  primitive.ObjectID `bson:"user_id"`  // For audience's queue
	PlaceID primitive.ObjectID `bson:"place_id"` // For place's queue
	Songs   []QueueSong        `bson:"songs"`
}

// SONG OFFERS (from artist to place)

type Status string

// Enum values for the status of a song offer.
const (
	Pending  Status = "pending"
	Accepted Status = "accepted"
	Rejected Status = "rejected"
)

// SongOffer represents an offer for a song.
type SongOffer struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ArtistID   primitive.ObjectID `bson:"artist_id"`
	PlaceID    primitive.ObjectID `bson:"place_id"`
	SongID     primitive.ObjectID `bson:"song_id"`
	Percentage float64            `bson:"percentage"`
	Status     Status             `bson:"status"`
}

type BulkSongOffer struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ArtistID   primitive.ObjectID `bson:"artist_id"`
	PlaceID    primitive.ObjectID `bson:"place_id"`
	Songs      []Music            `json:"songs,omitempty" bson:"songs,omitempty"`
	Percentage float64            `bson:"percentage"`
	Status     Status             `bson:"status"`
}
