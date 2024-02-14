package dao

import "go.mongodb.org/mongo-driver/bson/primitive"

type Music struct {
	ID     string             `json:"id,omitempty" bson:"_id,omitempty"`
	Place  primitive.ObjectID `json:"place"`
	Name   string             `json:"name"`
	Artist string             `json:"artist"`
	Path   string             `json:"path"`
}
