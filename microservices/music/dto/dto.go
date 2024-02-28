package dto

import "mime/multipart"

type MusicCreate struct {
	ID     string                `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string                `json:"name"`
	Artist string                `json:"artist"`
	Path   string                `json:"path"`
	Music  *multipart.FileHeader `json:"music"`
}

// MusicRetrieveDTO represents the DTO for retrieving an existing music record
type MusicRetrieve struct {
	ID     string `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string `json:"name,omitempty"`
	Artist string `json:"artist,omitempty"`
	Path   string `json:"path,omitempty"`
}
