package repository

import (
	"context"
	"errors"
	"reflect"
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
func (pr *PlaylistRepository) CreatePlaylist(playlist dao.Playlist) error {
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
func (pr *PlaylistRepository) GetPlaylistByID(id primitive.ObjectID) (dao.Playlist, error) {
	var playlist dao.Playlist

	filter := bson.M{"_id": id}

	err := collections.PlaylistsCollection.FindOne(context.Background(), filter).Decode(&playlist)
	if err != nil {
		return dao.Playlist{}, err
	}

	return playlist, nil
}

func (pr *PlaylistRepository) GetPlaylists() ([]dao.Playlist, error) {
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

// GetPlaylistsByOwner retrieves playlists by their owner's ID.
func (pr *PlaylistRepository) GetPlaylistsByOwner(ownerID primitive.ObjectID) ([]dao.Playlist, error) {
	var playlists []dao.Playlist

	filter := bson.M{"owner": ownerID}

	cursor, err := collections.PlaylistsCollection.Find(context.Background(), filter)
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
func (pr *PlaylistRepository) UpdatePlaylist(id primitive.ObjectID, playlist dao.Playlist) error {
	filter := bson.M{"_id": id}

	update := pr.generateUpdateQueryPlaylist(playlist)

	_, err := collections.PlaylistsCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// DeletePlaylist deletes a playlist from the database.
func (pr *PlaylistRepository) DeletePlaylist(id primitive.ObjectID) error {
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

// generateUpdateQueryPlaylist dynamically generates the update query based on the provided Playlist.
func (pr *PlaylistRepository) generateUpdateQueryPlaylist(playlist dao.Playlist) bson.M {
	update := bson.M{"$set": bson.M{}}

	playlistValue := reflect.ValueOf(playlist)
	playlistType := playlistValue.Type()

	for i := 0; i < playlistValue.NumField(); i++ {
		field := playlistValue.Field(i)
		fieldName := playlistType.Field(i).Name

		// Check if the field is a zero value or empty string
		if field.IsZero() && field.Interface() == "" {
			continue
		}

		update["$set"].(bson.M)[fieldName] = field.Interface()
	}

	// Set updatedAt field
	update["$set"].(bson.M)["updatedAt"] = time.Now()

	return update
}
