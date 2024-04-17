package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	dao "newnew.media/microservices/offer/dao"
)

type OfferRepository interface {
	CreateOffer(ctx context.Context, offer *dao.Offer) (*dao.Offer, error)
	GetOfferByID(ctx context.Context, id primitive.ObjectID) (*dao.Offer, error)
	GetOffersByPlace(ctx context.Context, place primitive.ObjectID) ([]*dao.Offer, error)
	GetOffersByArtist(ctx context.Context, artist primitive.ObjectID) ([]*dao.Offer, error)

	UpdateOfferStatus(ctx context.Context, id primitive.ObjectID, status dao.Status) error
	DeleteOffer(ctx context.Context, id primitive.ObjectID) error
}
