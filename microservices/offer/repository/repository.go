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
	offersCollection        *mongo.Collection
	counterOffersCollection *mongo.Collection
}

func NewMongoOfferRepository(db *mongo.Database) *OfferRepositoryImpl {
	return &OfferRepositoryImpl{
		offersCollection:        collections.OffersCollection,
		counterOffersCollection: collections.CounterOffersCollection,
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

func (r *OfferRepositoryImpl) UpdateOfferStatus(ctx context.Context, id primitive.ObjectID, status dao.Status) error {
	_, err := r.offersCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}

func (r *OfferRepositoryImpl) UpdateOfferWithCounterOffer(ctx context.Context, offerID primitive.ObjectID, counterOffer *dao.CounterOffer) error {
	offer, err := r.GetOfferByID(ctx, offerID)
	if err != nil {
		return err
	}

	offer.CounterOffers = append(offer.CounterOffers, *counterOffer)

	filter := bson.M{"_id": offerID}
	update := bson.M{"$set": bson.M{"counter_offers": offer.CounterOffers}}
	_, err = r.offersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *OfferRepositoryImpl) DeleteOffer(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.offersCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *OfferRepositoryImpl) CreateCounterOffer(ctx context.Context, counter *dao.Counter) (*dao.Counter, error) {
	result, err := r.counterOffersCollection.InsertOne(ctx, counter)
	if err != nil {
		return nil, err
	}
	counter.ID = result.InsertedID.(primitive.ObjectID)
	return counter, nil
}

func (r *OfferRepositoryImpl) GetCounterOfferByID(ctx context.Context, id primitive.ObjectID) (*dao.Counter, error) {
	var counter dao.Counter
	err := r.counterOffersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&counter)
	if err != nil {
		return nil, err
	}
	return &counter, nil
}

func (r *OfferRepositoryImpl) UpdateCounterOfferStatus(ctx context.Context, id primitive.ObjectID, status dao.Status) error {
	_, err := r.counterOffersCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}
