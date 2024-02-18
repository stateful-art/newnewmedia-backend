package dao

type Music struct {
	ID     string `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Path   string `json:"path"`
}
