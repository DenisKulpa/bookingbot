package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/DenisKulpa/bookingbot/internal/model"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// Create создаёт новое бронирование.
// clientID — внутренний id из таблицы users (ищем по telegram_id).
func (r *BookingRepository) Create(ctx context.Context, apartmentID, clientID int, checkIn, checkOut time.Time, guests int, totalPrice float64) (*model.Booking, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO bookings (apartment_id, client_id, check_in, check_out, guests_count, total_price, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'pending_approval')
		RETURNING id, apartment_id, client_id, check_in, check_out, guests_count, total_price, status,
		          COALESCE(admin_note,''), created_at, updated_at
	`, apartmentID, clientID, checkIn, checkOut, guests, totalPrice)

	b := &model.Booking{}
	if err := row.Scan(
		&b.ID, &b.ApartmentID, &b.ClientID,
		&b.CheckIn, &b.CheckOut,
		&b.GuestsCount, &b.TotalPrice, &b.Status,
		&b.AdminNote, &b.CreatedAt, &b.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("BookingRepository.Create: %w", err)
	}
	return b, nil
}

// GetBlockedDates возвращает все занятые даты для квартиры (из bookings + apartment_availability).
// Возвращает map[string]bool где ключ = "2006-01-02".
func (r *BookingRepository) GetBlockedDates(ctx context.Context, apartmentID int) (map[string]bool, error) {
	blocked := make(map[string]bool)

	// 1. Из подтверждённых/активных бронирований
	rows, err := r.db.QueryContext(ctx, `
		SELECT check_in, check_out FROM bookings
		WHERE apartment_id = $1
		  AND status IN ('approved', 'confirmed', 'payment_claimed')
		  AND check_out >= CURRENT_DATE
	`, apartmentID)
	if err != nil {
		return nil, fmt.Errorf("GetBlockedDates bookings: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var in, out time.Time
		if err := rows.Scan(&in, &out); err != nil {
			return nil, err
		}
		for d := in; d.Before(out); d = d.AddDate(0, 0, 1) {
			blocked[d.Format("2006-01-02")] = true
		}
	}

	// 2. Из apartment_availability (статус 'blocked')
	rows2, err := r.db.QueryContext(ctx, `
		SELECT date_from, date_to FROM apartment_availability
		WHERE apartment_id = $1
		  AND status = 'blocked'
		  AND date_to >= CURRENT_DATE
	`, apartmentID)
	if err != nil {
		return nil, fmt.Errorf("GetBlockedDates availability: %w", err)
	}
	defer rows2.Close()
	for rows2.Next() {
		var from, to time.Time
		if err := rows2.Scan(&from, &to); err != nil {
			return nil, err
		}
		for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
			blocked[d.Format("2006-01-02")] = true
		}
	}

	return blocked, nil
}

func (r *BookingRepository) GetOrCreateUser(ctx context.Context, telegramID int64, firstName, username string) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx,
		`SELECT id FROM users WHERE telegram_id = $1`, telegramID,
	).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return 0, fmt.Errorf("BookingRepository.GetOrCreateUser: %w", err)
	}
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO users (telegram_id, username, first_name, role)
		VALUES ($1, $2, $3, 'client')
		ON CONFLICT (telegram_id) DO UPDATE SET updated_at = NOW()
		RETURNING id
	`, telegramID, username, firstName).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("BookingRepository.GetOrCreateUser insert: %w", err)
	}
	return id, nil
}

// GetByID возвращает бронирование по id.
func (r *BookingRepository) GetByID(ctx context.Context, bookingID int) (*model.Booking, error) {
	b := &model.Booking{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, apartment_id, client_id, check_in, check_out,
		       guests_count, total_price, status, COALESCE(admin_note,''), created_at, updated_at
		FROM bookings WHERE id = $1
	`, bookingID).Scan(
		&b.ID, &b.ApartmentID, &b.ClientID,
		&b.CheckIn, &b.CheckOut,
		&b.GuestsCount, &b.TotalPrice, &b.Status,
		&b.AdminNote, &b.CreatedAt, &b.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("BookingRepository.GetByID: %w", err)
	}
	return b, nil
}

// UpdateStatus обновляет статус бронирования.
func (r *BookingRepository) UpdateStatus(ctx context.Context, bookingID int, status, note string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE bookings SET status = $2, admin_note = $3, updated_at = NOW()
		WHERE id = $1
	`, bookingID, status, note)
	return err
}

// GetOwnerTelegramID возвращает telegram_id владельца квартиры.
func (r *BookingRepository) GetOwnerTelegramID(ctx context.Context, apartmentID int) (int64, error) {
	var tgID int64
	err := r.db.QueryRowContext(ctx, `
		SELECT u.telegram_id FROM apartments a
		JOIN users u ON u.id = a.owner_id
		WHERE a.id = $1
	`, apartmentID).Scan(&tgID)
	if err != nil {
		return 0, fmt.Errorf("GetOwnerTelegramID: %w", err)
	}
	return tgID, nil
}

// GetClientTelegramID возвращает telegram_id клиента по его внутреннему id.
func (r *BookingRepository) GetClientTelegramID(ctx context.Context, clientID int) (int64, error) {
	var tgID int64
	err := r.db.QueryRowContext(ctx, `
		SELECT telegram_id FROM users WHERE id = $1
	`, clientID).Scan(&tgID)
	if err != nil {
		return 0, fmt.Errorf("GetClientTelegramID: %w", err)
	}
	return tgID, nil
}

// GetOwnerTelegramIDByBooking возвращает telegram_id владельца квартиры по id бронирования.
func (r *BookingRepository) GetOwnerTelegramIDByBooking(ctx context.Context, bookingID int) (int64, error) {
	var tgID int64
	err := r.db.QueryRowContext(ctx, `
		SELECT u.telegram_id
		FROM bookings b
		JOIN apartments a ON a.id = b.apartment_id
		JOIN users u ON u.id = a.owner_id
		WHERE b.id = $1
	`, bookingID).Scan(&tgID)
	if err != nil {
		return 0, fmt.Errorf("GetOwnerTelegramIDByBooking: %w", err)
	}
	return tgID, nil
}
