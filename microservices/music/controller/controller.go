package controller

import (
	"github.com/gofiber/fiber/v2"
)

// GetMusic gets all music
func GetMusic(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}
