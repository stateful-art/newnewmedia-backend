package dao

import "go.mongodb.org/mongo-driver/bson/primitive"

type Playlist struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Owner       string             `json:"owner,omitempty" bson:"owner,omitempty"`
}
