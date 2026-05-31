package model

import "time"

type ApartmentPhoto struct {
	ID          int       `json:"id"`
	ApartmentID int       `json:"apartment_id"`
	FilePath    string    `json:"file_path"` // uploads/apartments/{id}/filename.jpg
	URL         string    `json:"url"`        // /uploads/apartments/{id}/filename.jpg
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}
