package repository

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	collections "newnewmedia.com/db/collections"
	dao "newnewmedia.com/microservices/playlist/dao"
)

func GetPlaylists(c *fiber.Ctx) ([]dao.Playlist, error) {
	var playlists []dao.Playlist
	cursor, err := collections.PlaylistsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var playlist dao.Playlist
		cursor.Decode(&playlist)
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

func GetPlaylistById(c *fiber.Ctx, query *bson.M) (dao.Playlist, error) {
	var playlist dao.Playlist
	err := collections.PlaylistsCollection.FindOne(context.Background(), query).Decode(&playlist)
	if err != nil {
		return dao.Playlist{}, err
	}
	return playlist, nil
}
