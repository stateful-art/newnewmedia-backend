package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status string

const (
	Pending   Status = "pending"
	Accepted  Status = "accepted"
	Rejected  Status = "rejected"
	Countered Status = "countered"
)

type Preference string

const (
	Public     Preference = "public"
	Private    Preference = "private"
	Collective Preference = "collective"
	Individual Preference = "individual"
)

type Offer struct {
	ID          primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Songs       []primitive.ObjectID `json:"songs" bson:"songs"`
	Artist      primitive.ObjectID   `json:"artist" bson:"artist"`
	Place       primitive.ObjectID   `json:"place" bson:"place"`
	OfferedAt   time.Time            `json:"offered_at" bson:"offered_at"`
	ValidUntil  time.Time            `json:"valid_until" bson:"valid_until"`
	Status      Status               `json:"status" bson:"status"`
	Preferences []Preference         `json:"preferences" bson:"preferences"`

	// CounterOffers tracks the negotiation process.
	// Holds references to all counter offers made in response to the original offer.
	CounterOffers []CounterOffer `json:"counter_offers" bson:"counter_offers"`
}

type Counter struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Offer       primitive.ObjectID `json:"offer" bson:"offer"`
	OfferedAt   time.Time          `json:"offered_at" bson:"offered_at"`
	ValidUntil  time.Time          `json:"valid_until" bson:"valid_until"`
	Status      Status             `json:"status" bson:"status"`
	Preferences []Preference       `json:"preferences" bson:"preferences"`
	ParentOffer primitive.ObjectID `json:"parent_offer" bson:"parent_offer"`
}

// CounterOffer represents a counter offer (Counter) made by either the artist or the place.
// It will include a reference to the original offer and the counter offer itself.
type CounterOffer struct {
	CounterID primitive.ObjectID `json:"counter_id" bson:"counter_id"`
	Status    Status             `json:"status" bson:"status"`
}

// NegotiationHistory keeps track of all counter offers (CounterOffer) and their statuses
type NegotiationHistory struct {
	OriginalOfferID primitive.ObjectID `json:"original_offer_id" bson:"original_offer_id"`
	CounterOffers   []CounterOffer     `json:"counter_offers" bson:"counter_offers"`
}
