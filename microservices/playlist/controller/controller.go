package controller

import (
	"github.com/gofiber/fiber/v2"
	"newnewmedia.com/microservices/playlist/service"
)

func GetPlaylists(c *fiber.Ctx) error {
	playlists, err := service.GetPlaylists(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(playlists)
}
