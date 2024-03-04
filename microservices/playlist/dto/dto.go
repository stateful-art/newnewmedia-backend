package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Song struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Artist string             `json:"artist,omitempty" bson:"artist,omitempty"`
}

type Playlist struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Owner       string             `json:"owner,omitempty" bson:"owner,omitempty"`
	Songs       []Song             `json:"songs,omitempty" bson:"songs,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
