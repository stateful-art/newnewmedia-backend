package repository

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	collections "newnew.media/db/collections"
	dto "newnew.media/microservices/place/dto"
)

type PlaceRepository struct {
	// Any fields or dependencies needed by the repository can be added here
}

// NewPlaceRepository creates a new instance of the PlaceRepository.
func NewPlaceRepository() *PlaceRepository {
	return &PlaceRepository{}
}

func (pr *PlaceRepository) CreateGeospatialIndex() error {
	ctx := context.Background()
	indexModel := mongo.IndexModel{
		Keys: bson.M{"location": "2dsphere"},
	}
	_, err := collections.PlacesCollection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (pr *PlaceRepository) GetPlaces(c *fiber.Ctx) ([]dto.Place, error) {
	var places []dto.Place
	cursor, err := collections.PlacesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var place dto.Place
		cursor.Decode(&place)
		places = append(places, place)
	}
	return places, nil
}

func (pr *PlaceRepository) GetPlacesNearLocation(c *fiber.Ctx, longitude float64, latitude float64) ([]dto.Place, error) {
	var places []dto.Place

	// Define the geospatial query
	filter := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longitude, latitude},
				},
				"$maxDistance": 200, // The maximum distance in meters
			},
		},
	}

	cursor, err := collections.PlacesCollection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Error finding places: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var place dto.Place
		if err := cursor.Decode(&place); err != nil {
			log.Printf("Error decoding place: %v", err)
			continue // Skip this document if there's an error
		}
		places = append(places, place)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return places, nil
}

func (pr *PlaceRepository) CreatePlace(c *fiber.Ctx, place dto.Place) error {
	_, err := collections.PlacesCollection.InsertOne(context.Background(), place)
	if err != nil {
		return err
	}
	return nil
}

func (pr *PlaceRepository) GetPlace(c *fiber.Ctx, placeObjId *primitive.ObjectID) (dto.Place, error) {
	var place dto.Place
	err := collections.PlacesCollection.FindOne(context.Background(), bson.M{"_id": placeObjId}).Decode(&place)
	if err != nil {
		return dto.Place{}, err
	}
	return place, nil
}
