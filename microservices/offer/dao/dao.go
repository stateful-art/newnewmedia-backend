package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status string

const (
	Pending  Status = "pending"
	Accepted Status = "accepted"
	Rejected Status = "rejected"
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
}
