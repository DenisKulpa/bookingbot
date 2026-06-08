package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
		WHERE zone_id = $1
	`
	args := []any{zoneID}
	if onlyAvailable {
		q += " AND is_available = true"
	}
	q += " ORDER BY id ASC"
	return r.queryArgs(ctx, q, args...)
}

// GetByFilters возвращает квартиры, соответствующие фильтрам.
// Внутри одной категории — ИЛИ (любой из выбранных), между категориями — И (все категории).
// Если filterCodes пустой — возвращает все доступные квартиры.
func (r *ApartmentRepository) GetByFilters(ctx context.Context, filterCodes []string) ([]*model.Apartment, error) {
	if len(filterCodes) == 0 {
		return r.GetAllAvailable(ctx)
	}

	// Группируем коды по категориям (через БД)
	catMap, err := r.groupByCategory(ctx, filterCodes)
	if err != nil {
		return nil, err
	}

	// Строим условия: для каждой категории — EXISTS (хотя бы один код совпал)
	var conditions []string
	args := make([]any, 0)
	argIdx := 0
	for _, codes := range catMap {
		phs := make([]string, len(codes))
		for i, code := range codes {
			argIdx++
			phs[i] = fmt.Sprintf("$%d", argIdx)
			args = append(args, code)
		}
		conditions = append(conditions, fmt.Sprintf(`
			EXISTS (
				SELECT 1 FROM apartment_filters af
				JOIN filter_options fo ON fo.id = af.filter_option_id
				WHERE af.apartment_id = a.id
				  AND fo.code IN (%s)
			)`, strings.Join(phs, ",")))
	}

	q := fmt.Sprintf(`
		SELECT a.id, a.owner_id, a.zone_id, a.title, a.description, a.address,
		       a.rooms, a.max_guests, a.price_per_night, a.photos, a.amenities,
		       a.is_available, a.created_at, a.updated_at
		FROM apartments a
		WHERE a.is_available = true
		  AND %s
		ORDER BY a.id ASC
	`, strings.Join(conditions, " AND "))

	return r.queryArgs(ctx, q, args...)
}

// groupByCategory возвращает map[categoryCode][]filterCode
func (r *ApartmentRepository) groupByCategory(ctx context.Context, codes []string) (map[string][]string, error) {
	placeholders := make([]string, len(codes))
	args := make([]any, len(codes))
	for i, c := range codes {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = c
	}
	q := fmt.Sprintf(`
		SELECT DISTINCT fc.code, fo.code
		FROM filter_options fo
		JOIN filter_categories fc ON fc.id = fo.category_id
		WHERE fo.code IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]string)
	for rows.Next() {
		var catCode, optCode string
		if err := rows.Scan(&catCode, &optCode); err != nil {
			return nil, err
		}
		result[catCode] = append(result[catCode], optCode)
	}
	return result, rows.Err()
}

// GetAllAvailable возвращает все доступные квартиры (без привязки к конкретной зоне).
func (r *ApartmentRepository) GetAllAvailable(ctx context.Context) ([]*model.Apartment, error) {
	return r.queryArgs(ctx, `
		SELECT id, owner_id, zone_id, title, description, address,
		       rooms, max_guests, price_per_night, photos, amenities,
		       is_available, created_at, updated_at
		FROM apartments
		WHERE is_available = true
		ORDER BY id ASC
	`)
}

// GetAllAvailableLimited возвращает первые limit доступных квартир и общее количество.
func (r *ApartmentRepository) GetAllAvailableLimited(ctx context.Context, limit int) ([]*model.Apartment, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM apartments WHERE is_available = true`).Scan(&total); err != nil {
		return nil, 0, err
	}
	apts, err := r.queryArgs(ctx, `
		SELECT id, owner_id, zone_id, title, description, address,
		       rooms, max_guests, price_per_night, photos, amenities,
		       is_available, created_at, updated_at
		FROM apartments
		WHERE is_available = true
		ORDER BY id ASC
		LIMIT $1
	`, limit)
	return apts, total, err
}

func (r *ApartmentRepository) GetByID(ctx context.Context, id int) (*model.Apartment, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, owner_id, zone_id, title, description, address,
		       rooms, max_guests, price_per_night, photos, amenities,
		       is_available, created_at, updated_at
		FROM apartments
		WHERE id = $1
	`, id)

	a := &model.Apartment{}
	var (
		zoneID       sql.NullInt64
		description  sql.NullString
		address      sql.NullString
		photosRaw    string
		amenitiesRaw string
	)
	if err := row.Scan(
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

// Create создаёт новую квартиру.
func (r *ApartmentRepository) Create(ctx context.Context, ownerID int, zoneID *int, title, description, address, apartmentType string, rooms, maxGuests int, price float64) (*model.Apartment, error) {
	var zoneVal interface{}
	if zoneID != nil {
		zoneVal = *zoneID
	}
	var typeVal interface{}
	if apartmentType != "" {
		typeVal = apartmentType
	}
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO apartments (owner_id, zone_id, title, description, address, apartment_type, rooms, max_guests, price_per_night)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, owner_id, zone_id, title, description, address,
		          rooms, max_guests, price_per_night, photos, amenities,
		          is_available, created_at, updated_at
	`, ownerID, zoneVal, title, description, address, typeVal, rooms, maxGuests, price)

	a := &model.Apartment{}
	var (
		zID          sql.NullInt64
		desc         sql.NullString
		addr         sql.NullString
		photosRaw    string
		amenitiesRaw string
	)
	if err := row.Scan(
		&a.ID, &a.OwnerID, &zID,
		&a.Title, &desc, &addr,
		&a.Rooms, &a.MaxGuests, &a.PricePerNight,
		&photosRaw, &amenitiesRaw,
		&a.IsAvailable, &a.CreatedAt, &a.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("ApartmentRepository.Create: %w", err)
	}
	if zID.Valid {
		id := int(zID.Int64)
		a.ZoneID = &id
	}
	if desc.Valid {
		a.Description = desc.String
	}
	if addr.Valid {
		a.Address = addr.String
	}
	_ = a.ScanPhotos(photosRaw)
	_ = a.ScanAmenities(amenitiesRaw)
	return a, nil
}

// AddFilters привязывает filter_option коды к квартире.
func (r *ApartmentRepository) AddFilters(ctx context.Context, apartmentID int, filterCodes []string) error {
	if len(filterCodes) == 0 {
		return nil
	}
	for _, code := range filterCodes {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO apartment_filters (apartment_id, filter_option_id)
			SELECT $1, id FROM filter_options WHERE code = $2
			ON CONFLICT DO NOTHING
		`, apartmentID, code)
		if err != nil {
			return fmt.Errorf("ApartmentRepository.AddFilters (%s): %w", code, err)
		}
	}
	return nil
}

// GetByOwner возвращает все квартиры арендодателя по его internal user id.
func (r *ApartmentRepository) GetByOwner(ctx context.Context, ownerID int) ([]*model.Apartment, error) {
	return r.queryArgs(ctx, `
		SELECT id, owner_id, zone_id, title, description, address,
		       rooms, max_guests, price_per_night, photos, amenities,
		       is_available, created_at, updated_at
		FROM apartments WHERE owner_id = $1 ORDER BY id DESC
	`, ownerID)
}

// Update обновляет основные поля квартиры (title, description, address, price_per_night).
func (r *ApartmentRepository) Update(ctx context.Context, a *model.Apartment) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE apartments
		SET title = $1, description = $2, address = $3, price_per_night = $4, updated_at = NOW()
		WHERE id = $5
	`, a.Title, a.Description, a.Address, a.PricePerNight, a.ID)
	return err
}

// UpdateFull обновляет все редактируемые поля квартиры (включая zone, rooms, тип).
func (r *ApartmentRepository) UpdateFull(ctx context.Context, id int, zoneID *int, title, description, address, apartmentType string, rooms int, price float64) (int64, error) {
	var zoneVal interface{}
	if zoneID != nil {
		zoneVal = *zoneID
	}
	var typeVal interface{}
	if apartmentType != "" {
		typeVal = apartmentType
	}
	result, err := r.db.ExecContext(ctx, `
		UPDATE apartments
		SET zone_id = $1, title = $2, description = $3, address = $4,
		    apartment_type = $5, rooms = $6, max_guests = $7, price_per_night = $8,
		    updated_at = NOW()
		WHERE id = $9
	`, zoneVal, title, description, address, typeVal, rooms, rooms*2, price, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// ClearFilters удаляет все фильтры квартиры.
func (r *ApartmentRepository) ClearFilters(ctx context.Context, apartmentID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM apartment_filters WHERE apartment_id = $1`, apartmentID)
	return err
}

// GetFilterCodes возвращает все коды фильтров, привязанные к квартире.
func (r *ApartmentRepository) GetFilterCodes(ctx context.Context, apartmentID int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT fo.code FROM apartment_filters af
		JOIN filter_options fo ON fo.id = af.filter_option_id
		WHERE af.apartment_id = $1
	`, apartmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, rows.Err()
}

// GetAllCategories возвращает все активные категории фильтров с опциями.
func (r *ApartmentRepository) GetAllCategories(ctx context.Context) ([]model.FilterCategory, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT fc.code, fc.name, fo.code, fo.name
		FROM filter_categories fc
		JOIN filter_options fo ON fo.category_id = fc.id
		WHERE fc.is_active = 1 AND fo.is_active = 1
		ORDER BY fc.sort_order, fo.sort_order
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	catMap := make(map[string]*model.FilterCategory)
	var catOrder []string
	for rows.Next() {
		var catCode, catName, optCode, optName string
		if err := rows.Scan(&catCode, &catName, &optCode, &optName); err != nil {
			return nil, err
		}
		if _, ok := catMap[catCode]; !ok {
			catMap[catCode] = &model.FilterCategory{Code: catCode, Label: catName}
			catOrder = append(catOrder, catCode)
		}
		catMap[catCode].Options = append(catMap[catCode].Options, model.FilterOption{Code: optCode, Label: optName})
	}

	var result []model.FilterCategory
	for _, code := range catOrder {
		result = append(result, *catMap[code])
	}
	return result, rows.Err()
}
