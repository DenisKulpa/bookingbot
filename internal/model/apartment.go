package model

import (
	"encoding/json"
	"time"
)

type Apartment struct {
	ID           int       `json:"id"`
	OwnerID      int       `json:"owner_id"`
	ZoneID       *int      `json:"zone_id,omitempty"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	Address      string    `json:"address,omitempty"`
	Rooms        int       `json:"rooms"`
	MaxGuests    int       `json:"max_guests"`
	PricePerNight float64  `json:"price_per_night"`
	Photos       []string  `json:"photos"`
	Amenities    []string  `json:"amenities"`
	IsAvailable  bool      `json:"is_available"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (a *Apartment) PhotosJSON() (string, error) {
	b, err := json.Marshal(a.Photos)
	if err != nil {
		return "[]", err
	}
	return string(b), nil
}

func (a *Apartment) AmenitiesJSON() (string, error) {
	b, err := json.Marshal(a.Amenities)
	if err != nil {
		return "[]", err
	}
	return string(b), nil
}

func (a *Apartment) ScanPhotos(raw string) error {
	if raw == "" {
		a.Photos = []string{}
		return nil
	}
	return json.Unmarshal([]byte(raw), &a.Photos)
}

func (a *Apartment) ScanAmenities(raw string) error {
	if raw == "" {
		a.Amenities = []string{}
		return nil
	}
	return json.Unmarshal([]byte(raw), &a.Amenities)
}
