package repository

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	collections "newnewmedia.com/db/Collections"
	dao "newnewmedia.com/microservices/music/dao"
)

func CreateMusic(c *fiber.Ctx, music dao.Music) error {
	_, err := collections.MusicCollection.InsertOne(context.Background(), music)
	if err != nil {
		return err
	}
	return nil
}

func GetMusicByPlace(c *fiber.Ctx, id *primitive.ObjectID) ([]dao.Music, error) {
	var music []dao.Music
	cursor, err := collections.MusicCollection.Find(context.Background(), bson.M{"place": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var m dao.Music
		cursor.Decode(&m)
		music = append(music, m)
	}
	return music, nil
}
