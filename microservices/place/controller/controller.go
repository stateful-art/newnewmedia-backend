package controller

import (
	"github.com/gofiber/fiber/v2"
	"newnewmedia.com/microservices/place/dto"
	services "newnewmedia.com/microservices/place/service"
)

func GetPlaces(c *fiber.Ctx) error {
	places, err := services.GetPlaces(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Places Successfully Fetched",
		"data":    places,
	})
}

func CreatePlace(c *fiber.Ctx) error {
	var payload dto.Place
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	err := services.CreatePlace(c, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Place Successfully Created",
	})
}

func GetPlace(c *fiber.Ctx) error {
	id := c.Params("id")
	place, err := services.GetPlace(c, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Place Successfully Fetched",
		"data":    place,
	})
}
