package service

import (
	"context"

	dto "newnew.media/microservices/offer/dto"
)

type OfferService interface {
	CreateOffer(ctx context.Context, offer *dto.CreateOffer) (*dto.CreateOfferResponse, error)
	GetOfferByID(ctx context.Context, id string) (*dto.Offer, error)
	UpdateOfferStatus(ctx context.Context, id string, status string) error
	DeleteOffer(ctx context.Context, id string) error
	CreateCounterOffer(ctx context.Context, counter *dto.CreateCounterOffer) (*dto.CounterOffer, error)
	GetCounterOfferByID(ctx context.Context, id string) (*dto.Counter, error)
	UpdateCounterOfferStatus(ctx context.Context, id string, status string) error
}
