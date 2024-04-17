package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	collections "newnew.media/db/collections"
	dao "newnew.media/microservices/offer/dao"
)

type OfferRepositoryImpl struct {
	offersCollection *mongo.Collection
}

func NewMongoOfferRepository(db *mongo.Database) *OfferRepositoryImpl {
	return &OfferRepositoryImpl{
		offersCollection: collections.OffersCollection,
	}
}

func (r *OfferRepositoryImpl) CreateOffer(ctx context.Context, offer *dao.Offer) (*dao.Offer, error) {
	result, err := r.offersCollection.InsertOne(ctx, offer)
	if err != nil {
		return nil, err
	}
	offer.ID = result.InsertedID.(primitive.ObjectID)
	return offer, nil
}

func (r *OfferRepositoryImpl) GetOfferByID(ctx context.Context, id primitive.ObjectID) (*dao.Offer, error) {
	var offer dao.Offer
	err := r.offersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&offer)
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

func (r *OfferRepositoryImpl) GetOffersByPlace(ctx context.Context, place primitive.ObjectID) ([]*dao.Offer, error) {
	var offers []*dao.Offer // Change to slice of pointers
	cursor, err := r.offersCollection.Find(ctx, bson.M{"place": place})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each document into the offers slice
	for cursor.Next(ctx) {
		var offer dao.Offer
		if err := cursor.Decode(&offer); err != nil {
			return nil, err
		}
		offers = append(offers, &offer) // Append a pointer to the offer
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return offers, nil
}

func (r *OfferRepositoryImpl) GetOffersByArtist(ctx context.Context, artist primitive.ObjectID) ([]*dao.Offer, error) {
	var offers []*dao.Offer // Change to slice of pointers
	cursor, err := r.offersCollection.Find(ctx, bson.M{"artist": artist})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each document into the offers slice
	for cursor.Next(ctx) {
		var offer dao.Offer
		if err := cursor.Decode(&offer); err != nil {
			return nil, err
		}
		offers = append(offers, &offer) // Append a pointer to the offer
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return offers, nil
}

func (r *OfferRepositoryImpl) UpdateOfferStatus(ctx context.Context, id primitive.ObjectID, status dao.Status) error {
	_, err := r.offersCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}

func (r *OfferRepositoryImpl) DeleteOffer(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.offersCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
