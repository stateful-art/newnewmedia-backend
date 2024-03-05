package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	collections "newnewmedia.com/db/collections"
	dao "newnewmedia.com/microservices/revenue/dao"
)

type RevenueRepository struct {
	// Any fields or dependencies needed by the repository can be added here
}

// NewRevenueRepository creates a new instance of the RevenueRepository.
func NewRevenueRepository() *RevenueRepository {
	return &RevenueRepository{}
}

// CreateRevenue inserts a new revenue entry into the database.
func CreateRevenue(revenue dao.Revenue) error {
	revenue.ID = primitive.NewObjectID()
	revenue.CreatedAt = time.Now()
	revenue.UpdatedAt = time.Now()

	_, err := collections.RevenuesCollection.InsertOne(context.Background(), revenue)
	if err != nil {
		return err
	}
	return nil
}

// GetRevenueByID retrieves revenue by its ID.
func GetRevenueByID(id primitive.ObjectID) (dao.Revenue, error) {
	var revenue dao.Revenue

	filter := bson.M{"_id": id}

	err := collections.RevenuesCollection.FindOne(context.Background(), filter).Decode(&revenue)
	if err != nil {
		return dao.Revenue{}, err
	}

	return revenue, nil
}

// GetRevenueByArtistID retrieves revenue entries by artist ID.
func GetRevenueByArtistID(artistID primitive.ObjectID) ([]dao.Revenue, error) {
	var revenues []dao.Revenue

	filter := bson.M{"artist_id": artistID}

	cursor, err := collections.RevenuesCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var revenue dao.Revenue
		if err := cursor.Decode(&revenue); err != nil {
			return nil, err
		}
		revenues = append(revenues, revenue)
	}

	return revenues, nil
}

// GetRevenueByPlaceID retrieves revenue entries by place ID.
func GetRevenueByPlaceID(placeID primitive.ObjectID) ([]dao.Revenue, error) {
	var revenues []dao.Revenue

	filter := bson.M{"place_id": placeID}

	cursor, err := collections.RevenuesCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var revenue dao.Revenue
		if err := cursor.Decode(&revenue); err != nil {
			return nil, err
		}
		revenues = append(revenues, revenue)
	}

	return revenues, nil
}

func GetRevenueByPlaylistID(playlistID primitive.ObjectID) ([]dao.Revenue, error) {
	var revenues []dao.Revenue

	filter := bson.M{"playlist_id": playlistID}

	cursor, err := collections.RevenuesCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var revenue dao.Revenue
		if err := cursor.Decode(&revenue); err != nil {
			return nil, err
		}
		revenues = append(revenues, revenue)
	}

	return revenues, nil
}
