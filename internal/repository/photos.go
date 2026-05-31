package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DenisKulpa/bookingbot/internal/model"
)

type PhotoRepository struct {
	db *sql.DB
}

func NewPhotoRepository(db *sql.DB) *PhotoRepository {
	return &PhotoRepository{db: db}
}

// Add добавляет запись о фото после сохранения файла на диск.
func (r *PhotoRepository) Add(ctx context.Context, apartmentID int, filePath, url string, sortOrder int) (*model.ApartmentPhoto, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO apartment_photos (apartment_id, file_path, url, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, apartment_id, file_path, url, sort_order, created_at
	`, apartmentID, filePath, url, sortOrder)

	p := &model.ApartmentPhoto{}
	var urlNull sql.NullString
	if err := row.Scan(&p.ID, &p.ApartmentID, &p.FilePath, &urlNull, &p.SortOrder, &p.CreatedAt); err != nil {
		return nil, fmt.Errorf("PhotoRepository.Add scan: %w", err)
	}
	if urlNull.Valid {
		p.URL = urlNull.String
	}
	return p, nil
}

// GetByApartment возвращает все фото квартиры, отсортированные по sort_order.
func (r *PhotoRepository) GetByApartment(ctx context.Context, apartmentID int) ([]*model.ApartmentPhoto, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, apartment_id, file_path, url, sort_order, created_at
		FROM apartment_photos
		WHERE apartment_id = $1
		ORDER BY sort_order ASC, id ASC
	`, apartmentID)
	if err != nil {
		return nil, fmt.Errorf("PhotoRepository.GetByApartment: %w", err)
	}
	defer rows.Close()

	var photos []*model.ApartmentPhoto
	for rows.Next() {
		p := &model.ApartmentPhoto{}
		var urlNull sql.NullString
		if err := rows.Scan(&p.ID, &p.ApartmentID, &p.FilePath, &urlNull, &p.SortOrder, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("PhotoRepository.GetByApartment scan: %w", err)
		}
		if urlNull.Valid {
			p.URL = urlNull.String
		}
		photos = append(photos, p)
	}
	return photos, rows.Err()
}

// Delete удаляет запись о фото по id.
func (r *PhotoRepository) Delete(ctx context.Context, photoID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM apartment_photos WHERE id = $1`, photoID)
	return err
}
