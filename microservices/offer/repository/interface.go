package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	dao "newnew.media/microservices/offer/dao"
)

type OfferRepository interface {
	CreateOffer(ctx context.Context, offer *dao.Offer) (*dao.Offer, error)
	GetOfferByID(ctx context.Context, id primitive.ObjectID) (*dao.Offer, error)
	UpdateOfferStatus(ctx context.Context, id primitive.ObjectID, status dao.Status) error
	DeleteOffer(ctx context.Context, id primitive.ObjectID) error
	CreateCounterOffer(ctx context.Context, counter *dao.Counter) (*dao.Counter, error)
	GetCounterOfferByID(ctx context.Context, id primitive.ObjectID) (*dao.Counter, error)
	UpdateOfferWithCounterOffer(ctx context.Context, offerID primitive.ObjectID, counterOffer *dao.CounterOffer) error
	UpdateCounterOfferStatus(ctx context.Context, id primitive.ObjectID, status dao.Status) error
}
