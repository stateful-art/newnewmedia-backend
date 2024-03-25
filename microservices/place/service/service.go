package service

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"newnew.media/commons/utils"
	dto "newnew.media/microservices/place/dto"
	repository "newnew.media/microservices/place/repository"
)

type PlaceService struct {
	placeRepo  *repository.PlaceRepository
	natsClient *nats.Conn
}

func NewPlaceService(placeRepo *repository.PlaceRepository, natsClient *nats.Conn) *PlaceService {
	return &PlaceService{placeRepo: placeRepo, natsClient: natsClient}
}

func (ps *PlaceService) CreatePlace(c *fiber.Ctx, place dto.Place) error {
	err := ps.placeRepo.CreatePlace(c, place)
	if err != nil {
		return err
	}

	log.Print("marshaling place to index >>", place)
	// Serialize the place struct to JSON
	placeJSON, err := json.Marshal(place)
	if err != nil {
		log.Fatalf("Failed to serialize place to JSON: %v", err)
	}
	log.Print("marshalled place to index >>", placeJSON)

	// Send the JSON over NATS
	log.Print("Sending the JSON over NATS")
	erro := utils.SendMsgToPlaceIndexer(ps.natsClient, "place-created", placeJSON)
	if erro != nil {
		log.Print("Failed to index place", err)
	}
	log.Print("SENT the JSON over NATS")

	return nil
}

func (ps *PlaceService) GetPlace(c *fiber.Ctx, id string) (dto.Place, error) {
	placeObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return dto.Place{}, err
	}

	place, err := ps.placeRepo.GetPlace(c, &placeObjID)
	if err != nil {
		return dto.Place{}, err
	}
	return place, nil
}

func (ps *PlaceService) GetPlaces(c *fiber.Ctx) ([]dto.Place, error) {
	places, err := ps.placeRepo.GetPlaces(c)
	if err != nil {
		return nil, err
	}
	return places, nil
}
