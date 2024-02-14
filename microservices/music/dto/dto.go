package dto

import "mime/multipart"

type Music struct {
	ID     string                `json:"id,omitempty" bson:"_id,omitempty"`
	Place  string                `json:"place"`
	Name   string                `json:"name"`
	Artist string                `json:"artist"`
	Path   string                `json:"path"`
	Music  *multipart.FileHeader `json:"music"`
}
