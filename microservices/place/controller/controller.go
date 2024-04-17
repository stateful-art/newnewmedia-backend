package controller

import (
	"github.com/gofiber/fiber/v2"
	"newnew.media/microservices/place/dto"
	service "newnew.media/microservices/place/service"
)

type PlaceController struct {
	placeService *service.PlaceService
}

func NewPlaceController(placeService *service.PlaceService) *PlaceController {
	return &PlaceController{placeService: placeService}
}

func (pc *PlaceController) GetPlaces(c *fiber.Ctx) error {
	places, err := pc.placeService.GetPlaces(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(places)
}

func (pc *PlaceController) CreatePlace(c *fiber.Ctx) error {
	var payload dto.Place
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	err := pc.placeService.CreatePlace(c, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Place Successfully Created",
	})
}

func (pc *PlaceController) GetPlace(c *fiber.Ctx) error {
	id := c.Params("id")
	place, err := pc.placeService.GetPlace(c, id)
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
