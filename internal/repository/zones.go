package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/DenisKulpa/bookingbot/internal/model"
)

type ZoneRepository struct {
	db *sql.DB
}

func NewZoneRepository(db *sql.DB) *ZoneRepository {
	return &ZoneRepository{db: db}
}

func (r *ZoneRepository) GetTopLevel(ctx context.Context) ([]*model.Zone, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, emoji, short_desc, price_level, best_for, sort_order
		FROM zones
		WHERE (parent_id IS NULL OR parent_id = 0) AND is_active = 1
		ORDER BY sort_order ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("GetTopLevel query: %%w", err)
	}
	defer rows.Close()

	zones := make([]*model.Zone, 0)
	for rows.Next() {
		z := &model.Zone{}
		var emojiRaw, shortDescRaw, bestForRaw sql.NullString
		var priceLevelRaw sql.NullInt64
		if err := rows.Scan(
			&z.ID, &z.Name, &emojiRaw,
			&shortDescRaw, &priceLevelRaw, &bestForRaw, &z.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("GetTopLevel scan: %%w", err)
		}
		if emojiRaw.Valid {
			z.Emoji = emojiRaw.String
		}
		if shortDescRaw.Valid {
			z.ShortDesc = shortDescRaw.String
		}
		if bestForRaw.Valid {
			z.BestFor = bestForRaw.String
		}
		if priceLevelRaw.Valid {
			z.PriceLevel = int(priceLevelRaw.Int64)
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
	row := r.db.QueryRowContext(ctx, `
		SELECT
			id, parent_id, city, name, emoji,
			short_desc, full_desc, target_audience,
			pros, cons, housing_types,
			price_level, best_for, season_note, sort_order
		FROM zones
		WHERE id = ? AND is_active = 1
	`, id)

	z := &model.Zone{}
	var parentID sql.NullInt64
	var emojiRaw, shortDescRaw, fullDescRaw, targetAudienceRaw sql.NullString
	var prosRaw, consRaw, housingRaw sql.NullString
	var bestForRaw, seasonNoteRaw sql.NullString
	var priceLevelRaw sql.NullInt64

	err := row.Scan(
		&z.ID, &parentID, &z.City, &z.Name, &emojiRaw,
		&shortDescRaw, &fullDescRaw, &targetAudienceRaw,
		&prosRaw, &consRaw, &housingRaw,
		&priceLevelRaw, &bestForRaw, &seasonNoteRaw, &z.SortOrder,
	)
	if err != nil {
		return nil, fmt.Errorf("getZoneByID scan: %%w", err)
	}

	if emojiRaw.Valid {
		z.Emoji = emojiRaw.String
	}
	if shortDescRaw.Valid {
		z.ShortDesc = shortDescRaw.String
	}
	if fullDescRaw.Valid {
		z.FullDesc = fullDescRaw.String
	}
	if targetAudienceRaw.Valid {
		z.TargetAudience = targetAudienceRaw.String
	}
	if bestForRaw.Valid {
		z.BestFor = bestForRaw.String
	}
	if seasonNoteRaw.Valid {
		z.SeasonNote = seasonNoteRaw.String
	}
	if priceLevelRaw.Valid {
		z.PriceLevel = int(priceLevelRaw.Int64)
	}

	if parentID.Valid {
		v := int(parentID.Int64)
		z.ParentID = &v
	}
	if prosRaw.Valid {
		if prosRaw.String != "" {
			_ = json.Unmarshal([]byte(prosRaw.String), &z.Pros)
		}
	}
	if consRaw.Valid {
		if consRaw.String != "" {
			_ = json.Unmarshal([]byte(consRaw.String), &z.Cons)
		}
	}
	if housingRaw.Valid {
		if housingRaw.String != "" {
			_ = json.Unmarshal([]byte(housingRaw.String), &z.HousingTypes)
		}
	}

	return z, nil
}

func (r *ZoneRepository) getSubzones(ctx context.Context, parentID int) ([]*model.Zone, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, short_desc, full_desc, sort_order
		FROM zones
		WHERE parent_id = ? AND is_active = 1
		ORDER BY sort_order ASC
	`, parentID)
	if err != nil {
		return nil, fmt.Errorf("getSubzones query: %%w", err)
	}
	defer rows.Close()

	var subzones []*model.Zone
	for rows.Next() {
		z := &model.Zone{}
		var shortDescRaw, fullDescRaw sql.NullString
		if err := rows.Scan(&z.ID, &z.Name, &shortDescRaw, &fullDescRaw, &z.SortOrder); err != nil {
			return nil, fmt.Errorf("getSubzones scan: %%w", err)
		}
		if shortDescRaw.Valid {
			z.ShortDesc = shortDescRaw.String
		}
		if fullDescRaw.Valid {
			z.FullDesc = fullDescRaw.String
		}
		subzones = append(subzones, z)
	}

	return subzones, nil
}
