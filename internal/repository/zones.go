package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DenisKulpa/bookingbot/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ZoneRepository struct {
	db *pgxpool.Pool
}

func NewZoneRepository(db *pgxpool.Pool) *ZoneRepository {
	return &ZoneRepository{db: db}
}

func (r *ZoneRepository) GetTopLevel(ctx context.Context) ([]*model.Zone, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, emoji, short_desc, price_level, best_for, sort_order
		FROM zones
		WHERE parent_id IS NULL AND is_active = TRUE
		ORDER BY sort_order ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("GetTopLevel query: %w", err)
	}
	defer rows.Close()

	var zones []*model.Zone
	for rows.Next() {
		z := &model.Zone{}
		err := rows.Scan(
			&z.ID, &z.Name, &z.Emoji,
			&z.ShortDesc, &z.PriceLevel, &z.BestFor, &z.SortOrder,
		)
		if err != nil {
			return nil, fmt.Errorf("GetTopLevel scan: %w", err)
		}
		zones = append(zones, z)
	}

	return zones, nil
}

func (r *ZoneRepository) GetDistrictDetail(ctx context.Context, id int) (*model.DistrictDetail, error) {
	district, err := r.getZoneByID(ctx, id)
	if err != nil {
		return nil, err
	}

	subzones, err := r.getSubzones(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.DistrictDetail{
		District: district,
		Subzones: subzones,
	}, nil
}

func (r *ZoneRepository) getZoneByID(ctx context.Context, id int) (*model.Zone, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			id, parent_id, city, name, emoji,
			short_desc, full_desc, target_audience,
			pros, cons, housing_types,
			price_level, best_for, season_note, sort_order
		FROM zones
		WHERE id = $1 AND is_active = TRUE
	`, id)

	z := &model.Zone{}
	var prosRaw, consRaw, housingRaw []byte

	err := row.Scan(
		&z.ID, &z.ParentID, &z.City, &z.Name, &z.Emoji,
		&z.ShortDesc, &z.FullDesc, &z.TargetAudience,
		&prosRaw, &consRaw, &housingRaw,
		&z.PriceLevel, &z.BestFor, &z.SeasonNote, &z.SortOrder,
	)
	if err != nil {
		return nil, fmt.Errorf("getZoneByID scan: %w", err)
	}

	if prosRaw != nil {
		_ = json.Unmarshal(prosRaw, &z.Pros)
	}
	if consRaw != nil {
		_ = json.Unmarshal(consRaw, &z.Cons)
	}
	if housingRaw != nil {
		_ = json.Unmarshal(housingRaw, &z.HousingTypes)
	}

	return z, nil
}

func (r *ZoneRepository) getSubzones(ctx context.Context, parentID int) ([]*model.Zone, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, short_desc, full_desc, sort_order
		FROM zones
		WHERE parent_id = $1 AND is_active = TRUE
		ORDER BY sort_order ASC
	`, parentID)
	if err != nil {
		return nil, fmt.Errorf("getSubzones query: %w", err)
	}
	defer rows.Close()

	var subzones []*model.Zone
	for rows.Next() {
		z := &model.Zone{}
		if err := rows.Scan(&z.ID, &z.Name, &z.ShortDesc, &z.FullDesc, &z.SortOrder); err != nil {
			return nil, fmt.Errorf("getSubzones scan: %w", err)
		}
		subzones = append(subzones, z)
	}

	return subzones, nil
}
