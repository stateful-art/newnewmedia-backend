package service

import (
	"context"

	dto "newnew.media/microservices/offer/dto"
)

type OfferService interface {
	CreateOffer(ctx context.Context, offer *dto.CreateOffer) (*dto.CreateOfferResponse, error)
	GetOfferByID(ctx context.Context, id string) (*dto.Offer, error)
	GetOffersByPlace(ctx context.Context, id string) ([]*dto.Offer, error)
	GetOffersByArtist(ctx context.Context, id string) ([]*dto.Offer, error)

	UpdateOfferStatus(ctx context.Context, id string, status string) error
	DeleteOffer(ctx context.Context, id string) error
}
