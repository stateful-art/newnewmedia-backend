package service

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	dto "newnewmedia.com/microservices/place/dto"
	repository "newnewmedia.com/microservices/place/repository"
)

func CreatePlace(c *fiber.Ctx, place dto.Place) error {
	err := repository.CreatePlace(c, place)
	if err != nil {
		return err
	}
	return nil
}

func GetPlace(c *fiber.Ctx, id string) (dto.Place, error) {
	placeObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return dto.Place{}, err
	}

	place, err := repository.GetPlace(c, &placeObjID)
	if err != nil {
		return dto.Place{}, err
	}
	return place, nil
}

func GetPlaces(c *fiber.Ctx) ([]dto.Place, error) {
	places, err := repository.GetPlaces(c)
	if err != nil {
		return nil, err
	}
	return places, nil
}
