package repository

import (
	"context"

	"github.com/gofiber/fiber/v2"
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
