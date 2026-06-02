package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DenisKulpa/bookingbot/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

const userSelectCols = `
	id, telegram_id, username, first_name, last_name, role,
	is_blocked, COALESCE(phone,''), COALESCE(company_name,''), COALESCE(description,''),
	created_at, updated_at`

func scanUser(row interface{ Scan(...any) error }) (*model.User, error) {
	u := &model.User{}
	var isBlocked int
	err := row.Scan(
		&u.ID, &u.TelegramID, &u.Username, &u.FirstName, &u.LastName, &u.Role,
		&isBlocked, &u.Phone, &u.CompanyName, &u.Description,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	u.IsBlocked = isBlocked != 0
	return u, nil
}

// GetByID возвращает пользователя по внутреннему id.
func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+userSelectCols+` FROM users WHERE id = $1`, id)
	u, err := scanUser(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

// GetByTelegramID возвращает пользователя по telegram_id.
func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+userSelectCols+` FROM users WHERE telegram_id = $1`, telegramID)
	u, err := scanUser(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

// ListLandlords возвращает всех арендодателей.
func (r *UserRepository) ListLandlords(ctx context.Context) ([]*model.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+userSelectCols+` FROM users WHERE role = 'landlord' ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.ListLandlords: %w", err)
	}
	defer rows.Close()

	var result []*model.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("UserRepository.ListLandlords scan: %w", err)
		}
		result = append(result, u)
	}
	return result, rows.Err()
}

// CreateLandlord создаёт нового арендодателя.
func (r *UserRepository) CreateLandlord(ctx context.Context, telegramID int64, username, firstName, lastName, phone, companyName, description string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO users (telegram_id, username, first_name, last_name, role, phone, company_name, description)
		VALUES ($1, $2, $3, $4, 'landlord', $5, $6, $7)
		ON CONFLICT (telegram_id) DO UPDATE
		    SET role         = 'landlord',
		        username     = EXCLUDED.username,
		        first_name   = EXCLUDED.first_name,
		        last_name    = EXCLUDED.last_name,
		        phone        = EXCLUDED.phone,
		        company_name = EXCLUDED.company_name,
		        description  = EXCLUDED.description,
		        updated_at   = NOW()
		RETURNING `+userSelectCols,
		telegramID, username, firstName, lastName, phone, companyName, description,
	)
	u, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.CreateLandlord: %w", err)
	}
	return u, nil
}

// UpdateLandlord обновляет профиль арендодателя.
func (r *UserRepository) UpdateLandlord(ctx context.Context, id int, phone, companyName, description string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE users
		SET phone = $2, company_name = $3, description = $4, updated_at = NOW()
		WHERE id = $1 AND role = 'landlord'
		RETURNING `+userSelectCols,
		id, phone, companyName, description,
	)
	u, err := scanUser(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("UserRepository.UpdateLandlord: %w", err)
	}
	return u, nil
}

// DeleteLandlord удаляет арендодателя (и каскадно все его квартиры).
func (r *UserRepository) DeleteLandlord(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM users WHERE id = $1 AND role = 'landlord'`, id)
	return err
}

// ListApartmentsByLandlord возвращает все квартиры арендодателя.
func (r *UserRepository) ListApartmentsByLandlord(ctx context.Context, landlordID int) ([]map[string]interface{}, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT a.id, a.title, a.address, a.price_per_night, a.is_available, a.created_at
		FROM apartments a
		WHERE a.owner_id = $1
		ORDER BY a.id`, landlordID)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.ListApartmentsByLandlord: %w", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id int
		var title, address string
		var price float64
		var isAvailable bool
		var createdAt interface{}
		if err := rows.Scan(&id, &title, &address, &price, &isAvailable, &createdAt); err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{
			"id":            id,
			"title":         title,
			"address":       address,
			"price_per_night": price,
			"is_available":  isAvailable,
			"created_at":    createdAt,
		})
	}
	return result, rows.Err()
}
