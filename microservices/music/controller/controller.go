package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	dto "newnewmedia.com/microservices/music/dto"
	service "newnewmedia.com/microservices/music/service"
)

// GetMusic gets all music
func GetMusic(c *fiber.Ctx) error {

	return c.SendFile("./public/music/test.mp3")
}

func GetMusicFile(c *fiber.Ctx) error {
	fileName := c.Params("id")
	return c.SendFile("./public/music/" + fileName)
}

func CreateMusic(c *fiber.Ctx) error {
	var musicPayload dto.Music
	if err := c.BodyParser(&musicPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if musicPayload.Name == "" || musicPayload.Artist == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Name and Artist are required",
		})
	}
	log.Println(musicPayload)
	err := service.CreateMusic(c, musicPayload)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Music created successfully",
	})

}

func GetMusicByPlace(c *fiber.Ctx) error {
	id := c.Params("id")
	music, err := service.GetMusicByPlace(c, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Music fetched successfully",
		"data":    music,
	})
}
