package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DenisKulpa/bookingbot/internal/model"
)

type ApartmentRepository struct {
	db *sql.DB
}

func NewApartmentRepository(db *sql.DB) *ApartmentRepository {
	return &ApartmentRepository{db: db}
}

func (r *ApartmentRepository) GetByZone(ctx context.Context, zoneID int, onlyAvailable bool) ([]*model.Apartment, error) {
	q := `
		SELECT id, owner_id, zone_id, title, description, address,
		       rooms, max_guests, price_per_night, photos, amenities,
		       is_available, created_at, updated_at
		FROM apartments
		WHERE zone_id = ?
	`
	args := []any{zoneID}
	if onlyAvailable {
		q += " AND is_available = true"
	}
	q += " ORDER BY id ASC"
	return r.queryArgs(ctx, q, args...)
}

func (r *ApartmentRepository) queryArgs(ctx context.Context, q string, args ...any) ([]*model.Apartment, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("apartments query: %w", err)
	}
	defer rows.Close()

	apartments := make([]*model.Apartment, 0)
	for rows.Next() {
		a, err := scanApartment(rows)
		if err != nil {
			return nil, fmt.Errorf("apartments scan: %w", err)
		}
		apartments = append(apartments, a)
	}

	return apartments, rows.Err()
}

func scanApartment(rows *sql.Rows) (*model.Apartment, error) {
	a := &model.Apartment{}
	var (
		zoneID       sql.NullInt64
		description  sql.NullString
		address      sql.NullString
		photosRaw    string
		amenitiesRaw string
	)

	if err := rows.Scan(
		&a.ID, &a.OwnerID, &zoneID,
		&a.Title, &description, &address,
		&a.Rooms, &a.MaxGuests, &a.PricePerNight,
		&photosRaw, &amenitiesRaw,
		&a.IsAvailable, &a.CreatedAt, &a.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if zoneID.Valid {
		id := int(zoneID.Int64)
		a.ZoneID = &id
	}
	if description.Valid {
		a.Description = description.String
	}
	if address.Valid {
		a.Address = address.String
	}
	if err := a.ScanPhotos(photosRaw); err != nil {
		return nil, fmt.Errorf("parse photos: %w", err)
	}
	if err := a.ScanAmenities(amenitiesRaw); err != nil {
		return nil, fmt.Errorf("parse amenities: %w", err)
	}

	return a, nil
}
