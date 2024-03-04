package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	collections "newnewmedia.com/db/collections"
	dao "newnewmedia.com/microservices/playlist/dao"
)

type PlaylistRepository struct {
	// Any fields or dependencies needed by the repository can be added here
}

// NewPlaylistRepository creates a new instance of the PlaylistRepository.
func NewPlaylistRepository() *PlaylistRepository {
	return &PlaylistRepository{}
}

// CreatePlaylist inserts a new playlist into the database.
func CreatePlaylist(playlist dao.Playlist) error {
	playlist.ID = primitive.NewObjectID()
	playlist.CreatedAt = time.Now()
	playlist.UpdatedAt = time.Now()

	_, err := collections.PlaylistsCollection.InsertOne(context.Background(), playlist)
	if err != nil {
		return err
	}
	return nil
}

// GetPlaylistByID retrieves a playlist by its ID.
func GetPlaylistByID(id primitive.ObjectID) (dao.Playlist, error) {
	var playlist dao.Playlist

	filter := bson.M{"_id": id}

	err := collections.PlaylistsCollection.FindOne(context.Background(), filter).Decode(&playlist)
	if err != nil {
		return dao.Playlist{}, err
	}

	return playlist, nil
}

func GetPlaylists() ([]dao.Playlist, error) {
	var playlists []dao.Playlist

	cursor, err := collections.PlaylistsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var playlist dao.Playlist
		if err := cursor.Decode(&playlist); err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

// UpdatePlaylist updates an existing playlist in the database.
func UpdatePlaylist(id primitive.ObjectID, playlist dao.Playlist) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"name":        playlist.Name,
			"description": playlist.Description,
			"owner":       playlist.Owner,
			"updatedAt":   time.Now(),
			"songs":       playlist.Songs,
		},
	}

	_, err := collections.PlaylistsCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// DeletePlaylist deletes a playlist from the database.
func DeletePlaylist(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := collections.PlaylistsCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no document found to delete")
	}

	return nil
}
