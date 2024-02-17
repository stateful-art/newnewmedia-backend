package service

import (
	"github.com/gofiber/fiber/v2"
	dao "newnewmedia.com/microservices/playlist/dao"
	repository "newnewmedia.com/microservices/playlist/repository"
)

func GetPlaylists(c *fiber.Ctx) ([]dao.Playlist, error) {
	playlists, err := repository.GetPlaylists(c)
	if err != nil {
		return nil, err
	}
	return playlists, nil
}
