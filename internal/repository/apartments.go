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

func (r *ApartmentRepository) GetByZone(ctx context.Context, SubzoneID int, onlyAvailable bool) ([]*model.Apartment, error) {
	q := `
		SELECT id, owner_id, subzone_id, title, description, address,
		       rooms, max_guests, price_per_night, photos, amenities,
		       is_available, created_at, updated_at
		FROM apartments
		WHERE subzone_id = $1
	`
	args := []any{SubzoneID}
	if onlyAvailable {
		q += " AND is_available = true"
	}
	q += " ORDER BY id ASC"
	return r.queryArgs(ctx, q, args...)
}

// GetByFilters возвращает квартиры по фильтрам в указанном городе.
// city="" — без фильтрации по городу.
func (r *ApartmentRepository) GetByFilters(ctx context.Context, city string, filterCodes []string) ([]*model.Apartment, error) {
	cityID, err := r.resolveCityID(ctx, city)
	if err != nil {
		return nil, err
	}
	if len(filterCodes) == 0 {
		return r.GetAllAvailable(ctx, cityID)
	}

	catMap, err := r.groupByCategory(ctx, filterCodes)
	if err != nil {
		return nil, err
	}

	var conditions []string
	args := make([]any, 0)
	argIdx := 0

	if cityID > 0 {
		argIdx++
		conditions = append(conditions, fmt.Sprintf(`a.city_id = $%d`, argIdx))
		args = append(args, cityID)
	}

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
				WHERE af.apartment_id = a.id AND fo.code IN (%s)
			)`, strings.Join(phs, ",")))
	}

	q := fmt.Sprintf(`
		SELECT a.id, a.owner_id, a.subzone_id, a.city_id, a.title, a.description, a.address,
		       a.rooms, a.max_guests, a.price_per_night, a.photos, a.amenities,
		       a.is_available, a.created_at, a.updated_at
		FROM apartments a
		WHERE a.is_available = true AND %s
		ORDER BY a.id ASC
	`, strings.Join(conditions, " AND "))

	return r.queryArgs(ctx, q, args...)
}

func (r *ApartmentRepository) resolveCityID(ctx context.Context, city string) (int, error) {
	if city == "" {
		return 0, nil
	}
	var id int
	err := r.db.QueryRowContext(ctx, `SELECT id FROM cities WHERE name = $1`, city).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}

// GetAllAvailable возвращает доступные квартиры в городе (cityID=0 — все).
func (r *ApartmentRepository) GetAllAvailable(ctx context.Context, cityID int) ([]*model.Apartment, error) {
	return r.queryAvailable(ctx, cityID, 0)
}

func (r *ApartmentRepository) GetAllAvailableLimited(ctx context.Context, city string, limit int) ([]*model.Apartment, int, error) {
	cityID, err := r.resolveCityID(ctx, city)
	if err != nil {
		return nil, 0, err
	}
	return r.getAllAvailableLimited(ctx, cityID, limit)
}

func (r *ApartmentRepository) getAllAvailableLimited(ctx context.Context, cityID int, limit int) ([]*model.Apartment, int, error) {
	where := "WHERE a.is_available = true"
	args := make([]any, 0)
	argIdx := 1
	if cityID > 0 {
		where += fmt.Sprintf(` AND a.city_id = $%d`, argIdx)
		args = append(args, cityID)
		argIdx++
	}
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM apartments a `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, limit)
	apts, err := r.queryArgs(ctx, fmt.Sprintf(`
		SELECT a.id, a.owner_id, a.subzone_id, a.city_id, a.title, a.description, a.address,
		       a.rooms, a.max_guests, a.price_per_night, a.photos, a.amenities,
		       a.is_available, a.created_at, a.updated_at
		FROM apartments a
		%s
		ORDER BY a.id ASC
		LIMIT $%d
	`, where, argIdx), args...)
	return apts, total, err
}

func (r *ApartmentRepository) queryAvailable(ctx context.Context, cityID int, limit int) ([]*model.Apartment, error) {
	where := "WHERE a.is_available = true"
	args := make([]any, 0)
	argIdx := 1
	if cityID > 0 {
		where += fmt.Sprintf(` AND a.city_id = $%d`, argIdx)
		args = append(args, cityID)
		argIdx++
	}
	q := fmt.Sprintf(`
		SELECT a.id, a.owner_id, a.subzone_id, a.city_id, a.title, a.description, a.address,
		       a.rooms, a.max_guests, a.price_per_night, a.photos, a.amenities,
		       a.is_available, a.created_at, a.updated_at
		FROM apartments a
		%s
		ORDER BY a.id ASC
	`, where)
	if limit > 0 {
		q += fmt.Sprintf(` LIMIT $%d`, argIdx)
		args = append(args, limit)
	}
	return r.queryArgs(ctx, q, args...)
}

func (r *ApartmentRepository) GetByID(ctx context.Context, id int) (*model.Apartment, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, owner_id, subzone_id, title, description, address,
		       rooms, max_guests, price_per_night, photos, amenities,
		       is_available, created_at, updated_at
		FROM apartments
		WHERE id = $1
	`, id)

	a := &model.Apartment{}
	var (
		SubzoneID    sql.NullInt64
		cityID       sql.NullInt64
		description  sql.NullString
		address      sql.NullString
		photosRaw    string
		amenitiesRaw string
	)
	if err := row.Scan(
		&a.ID, &a.OwnerID, &SubzoneID, &cityID,
		&a.Title, &description, &address,
		&a.Rooms, &a.MaxGuests, &a.PricePerNight,
		&photosRaw, &amenitiesRaw,
		&a.IsAvailable, &a.CreatedAt, &a.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if SubzoneID.Valid {
		id := int(SubzoneID.Int64)
		a.SubzoneID = &id
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
		SubzoneID    sql.NullInt64
		cityID       sql.NullInt64
		description  sql.NullString
		address      sql.NullString
		photosRaw    string
		amenitiesRaw string
	)

	if err := rows.Scan(
		&a.ID, &a.OwnerID, &SubzoneID, &cityID,
		&a.Title, &description, &address,
		&a.Rooms, &a.MaxGuests, &a.PricePerNight,
		&photosRaw, &amenitiesRaw,
		&a.IsAvailable, &a.CreatedAt, &a.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if SubzoneID.Valid {
		id := int(SubzoneID.Int64)
		a.SubzoneID = &id
	}
	if cityID.Valid {
		id := int(cityID.Int64)
		a.CityID = &id
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
func (r *ApartmentRepository) Create(ctx context.Context, ownerID int, SubzoneID *int, title, description, address, apartmentType string, rooms, maxGuests int, price float64) (*model.Apartment, error) {
	var zoneVal interface{}
	if SubzoneID != nil {
		zoneVal = *SubzoneID
	}
	var typeVal interface{}
	if apartmentType != "" {
		typeVal = apartmentType
	}
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO apartments (owner_id, subzone_id, city_id, title, description, address, apartment_type, rooms, max_guests, price_per_night)
		SELECT $1, $2, c.id, $3, $4, $5, $6, $7, $8, $9
		FROM subzones s JOIN zones z ON z.id = s.zone_id JOIN cities c ON c.id = z.city_id
		WHERE s.id = $2
		RETURNING id, owner_id, subzone_id, city_id, title, description, address,
		          rooms, max_guests, price_per_night, photos, amenities,
		          is_available, created_at, updated_at
	`, ownerID, zoneVal, title, description, address, typeVal, rooms, maxGuests, price)

	a := &model.Apartment{}
	var (
		zID          sql.NullInt64
		cID          sql.NullInt64
		desc         sql.NullString
		addr         sql.NullString
		photosRaw    string
		amenitiesRaw string
	)
	if err := row.Scan(
		&a.ID, &a.OwnerID, &zID, &cID,
		&a.Title, &desc, &addr,
		&a.Rooms, &a.MaxGuests, &a.PricePerNight,
		&photosRaw, &amenitiesRaw,
		&a.IsAvailable, &a.CreatedAt, &a.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("ApartmentRepository.Create: %w", err)
	}
	if zID.Valid {
		id := int(zID.Int64)
		a.SubzoneID = &id
	}
	if cID.Valid {
		id := int(cID.Int64)
		a.CityID = &id
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
		SELECT id, owner_id, subzone_id, title, description, address,
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
func (r *ApartmentRepository) UpdateFull(ctx context.Context, id int, SubzoneID *int, title, description, address, apartmentType string, rooms int, price float64) (int64, error) {
	var zoneVal interface{}
	if SubzoneID != nil {
		zoneVal = *SubzoneID
	}
	var typeVal interface{}
	if apartmentType != "" {
		typeVal = apartmentType
	}
	result, err := r.db.ExecContext(ctx, `
		UPDATE apartments
		SET subzone_id = $1, city_id = (SELECT c.id FROM subzones s JOIN zones z ON z.id = s.zone_id JOIN cities c ON c.id = z.city_id WHERE s.id = $1),
		    title = $2, description = $3, address = $4,
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



func (r *ApartmentRepository) groupByCategory(ctx context.Context, codes []string) (map[string][]string, error) {
	phs := make([]string, len(codes))
	args := make([]any, len(codes))
	for i, c := range codes {
		phs[i] = fmt.Sprintf("$%d", i+1)
		args[i] = c
	}
	q := fmt.Sprintf("SELECT DISTINCT fc.code, fo.code FROM filter_options fo JOIN filter_categories fc ON fc.id = fo.category_id WHERE fo.code IN (%s)", strings.Join(phs, ","))
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	result := make(map[string][]string)
	for rows.Next() {
		var catCode, optCode string
		if err := rows.Scan(&catCode, &optCode); err != nil { return nil, err }
		result[catCode] = append(result[catCode], optCode)
	}
	return result, rows.Err()
}
