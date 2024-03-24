package dto

import (
	"mime/multipart"
)

type CreateMusic struct {
	Name   string                `json:"name"`
	Artist string                `json:"artist"`
	Path   string                `json:"path"`
	Music  *multipart.FileHeader `json:"music"`
	Genres []string              `json:"genres"`
}

// MusicRetrieveDTO represents the DTO for retrieving an existing music record
type MusicRetrieve struct {
	ID     string   `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string   `json:"name,omitempty"`
	Artist string   `json:"artist,omitempty"`
	Path   string   `json:"path,omitempty"`
	Genres []string `json:"genres"`
}
