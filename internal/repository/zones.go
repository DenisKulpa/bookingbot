package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DenisKulpa/bookingbot/internal/model"
)

type ZoneRepository struct {
	db *sql.DB
}

func NewZoneRepository(db *sql.DB) *ZoneRepository {
	return &ZoneRepository{db: db}
}

// ─── Cities ───────────────────────────────────────────────────────────────────

func (r *ZoneRepository) GetCities(ctx context.Context) ([]*model.City, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, emoji, sort_order FROM cities WHERE is_active = 1 ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cities []*model.City
	for rows.Next() {
		c := &model.City{}
		var emoji sql.NullString
		if err := rows.Scan(&c.ID, &c.Name, &emoji, &c.SortOrder); err != nil {
			return nil, err
		}
		if emoji.Valid {
			c.Emoji = emoji.String
		}
		cities = append(cities, c)
	}
	return cities, rows.Err()
}

// ─── Zones ────────────────────────────────────────────────────────────────────

func (r *ZoneRepository) GetZonesByCity(ctx context.Context, cityID int) ([]*model.Zone, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, emoji, short_desc, price_level, best_for, sort_order FROM zones WHERE city_id = $1 AND is_active = 1 ORDER BY sort_order`, cityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanZones(rows)
}

func scanZones(rows *sql.Rows) ([]*model.Zone, error) {
	var zones []*model.Zone
	for rows.Next() {
		z := &model.Zone{}
		var emoji, shortDesc, bestFor sql.NullString
		var priceLevel sql.NullInt64
		if err := rows.Scan(&z.ID, &z.Name, &emoji, &shortDesc, &priceLevel, &bestFor, &z.SortOrder); err != nil {
			return nil, err
		}
		if emoji.Valid { z.Emoji = emoji.String }
		if shortDesc.Valid { z.ShortDesc = shortDesc.String }
		if bestFor.Valid { z.BestFor = bestFor.String }
		if priceLevel.Valid { z.PriceLevel = int(priceLevel.Int64) }
		zones = append(zones, z)
	}
	return zones, rows.Err()
}

// ─── Subzones ─────────────────────────────────────────────────────────────────

func (r *ZoneRepository) GetSubzonesByZone(ctx context.Context, zoneID int) ([]*model.Subzone, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, emoji, short_desc, sort_order FROM subzones WHERE zone_id = $1 AND is_active = 1 ORDER BY sort_order`, zoneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subzones []*model.Subzone
	for rows.Next() {
		s := &model.Subzone{ZoneID: zoneID}
		var emoji, shortDesc sql.NullString
		if err := rows.Scan(&s.ID, &s.Name, &emoji, &shortDesc, &s.SortOrder); err != nil {
			return nil, err
		}
		if emoji.Valid { s.Emoji = emoji.String }
		if shortDesc.Valid { s.ShortDesc = shortDesc.String }
		subzones = append(subzones, s)
	}
	return subzones, rows.Err()
}

func (r *ZoneRepository) GetAssignableSubzones(ctx context.Context) ([]*model.Subzone, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.zone_id, s.name, s.emoji, s.short_desc, s.sort_order, z.name, c.name
		FROM subzones s
		JOIN zones z ON z.id = s.zone_id
		JOIN cities c ON c.id = z.city_id
		WHERE s.is_active = 1 AND z.is_active = 1 AND c.is_active = 1
		ORDER BY c.sort_order, z.sort_order, s.sort_order
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subzones []*model.Subzone
	for rows.Next() {
		s := &model.Subzone{}
		var emoji, shortDesc sql.NullString
		if err := rows.Scan(&s.ID, &s.ZoneID, &s.Name, &emoji, &shortDesc, &s.SortOrder, &s.ZoneName, &s.CityName); err != nil {
			return nil, err
		}
		if emoji.Valid { s.Emoji = emoji.String }
		if shortDesc.Valid { s.ShortDesc = shortDesc.String }
		subzones = append(subzones, s)
	}
	return subzones, rows.Err()
}

func (r *ZoneRepository) GetSubzoneTree(ctx context.Context) ([]*model.City, map[int][]*model.Zone, map[int][]*model.Subzone, error) {
	cities, err := r.GetCities(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	zRows, err := r.db.QueryContext(ctx, `SELECT id, city_id, name, emoji, short_desc, sort_order FROM zones WHERE is_active = 1 ORDER BY sort_order`)
	if err != nil {
		return nil, nil, nil, err
	}
	defer zRows.Close()
	zoneMap := make(map[int][]*model.Zone)
	for zRows.Next() {
		z := &model.Zone{}
		var emoji, shortDesc sql.NullString
		if err := zRows.Scan(&z.ID, &z.CityID, &z.Name, &emoji, &shortDesc, &z.SortOrder); err != nil {
			return nil, nil, nil, err
		}
		if emoji.Valid { z.Emoji = emoji.String }
		if shortDesc.Valid { z.ShortDesc = shortDesc.String }
		zoneMap[z.CityID] = append(zoneMap[z.CityID], z)
	}

	sRows, err := r.db.QueryContext(ctx, `SELECT id, zone_id, name, emoji, short_desc, sort_order FROM subzones WHERE is_active = 1 ORDER BY sort_order`)
	if err != nil {
		return nil, nil, nil, err
	}
	defer sRows.Close()
	subMap := make(map[int][]*model.Subzone)
	for sRows.Next() {
		s := &model.Subzone{}
		var emoji, shortDesc sql.NullString
		if err := sRows.Scan(&s.ID, &s.ZoneID, &s.Name, &emoji, &shortDesc, &s.SortOrder); err != nil {
			return nil, nil, nil, err
		}
		if emoji.Valid { s.Emoji = emoji.String }
		if shortDesc.Valid { s.ShortDesc = shortDesc.String }
		subMap[s.ZoneID] = append(subMap[s.ZoneID], s)
	}
	return cities, zoneMap, subMap, nil
}

func (r *ZoneRepository) GetFilterCodes(ctx context.Context, subzoneIDs ...int) ([]string, error) {
	if len(subzoneIDs) == 0 {
		return nil, nil
	}
	phs := make([]string, len(subzoneIDs))
	args := make([]any, len(subzoneIDs))
	for i, id := range subzoneIDs {
		phs[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	q := fmt.Sprintf(`SELECT DISTINCT fo.code FROM filter_options fo JOIN subzones s ON s.name = fo.name WHERE s.id IN (%s) AND fo.code LIKE 'zone_%%'`, strings.Join(phs, ","))
	rows, err := r.db.QueryContext(ctx, q, args...)
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

func (r *ZoneRepository) GetCityFilterCodes(ctx context.Context, cityName string) (map[string]bool, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT DISTINCT fo.code FROM filter_options fo
		JOIN subzones s ON s.name = fo.name
		JOIN zones z ON z.id = s.zone_id
		JOIN cities c ON c.id = z.city_id
		WHERE c.name = $1 AND fo.code LIKE 'zone_%'
	`, cityName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]bool)
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		result[code] = true
	}
	return result, rows.Err()
}

// ─── Legacy (для API) ────────────────────────────────────────────────────────

func (r *ZoneRepository) GetTopLevel(ctx context.Context) ([]*model.Zone, error) {
	return r.GetZonesByCity(ctx, 1)
}

func (r *ZoneRepository) GetDistrictDetail(ctx context.Context, id int) (*model.DistrictDetail, error) {
	z := &model.Zone{}
	row := r.db.QueryRowContext(ctx, `SELECT id, city_id, name, emoji, short_desc, full_desc, target_audience, pros, cons, housing_types, price_level, best_for, season_note, sort_order FROM zones WHERE id = $1 AND is_active = 1`, id)
	var cityID int
	var emoji, shortDesc, fullDesc, targetAud, bestFor, season sql.NullString
	var prosRaw, consRaw, housingRaw sql.NullString
	var priceLevel sql.NullInt64
	if err := row.Scan(&z.ID, &cityID, &z.Name, &emoji, &shortDesc, &fullDesc, &targetAud, &prosRaw, &consRaw, &housingRaw, &priceLevel, &bestFor, &season, &z.SortOrder); err != nil {
		return nil, err
	}
	if emoji.Valid { z.Emoji = emoji.String }
	if shortDesc.Valid { z.ShortDesc = shortDesc.String }
	if fullDesc.Valid { z.FullDesc = fullDesc.String }
	if targetAud.Valid { z.TargetAudience = targetAud.String }
	if bestFor.Valid { z.BestFor = bestFor.String }
	if season.Valid { z.SeasonNote = season.String }
	if priceLevel.Valid { z.PriceLevel = int(priceLevel.Int64) }
	if prosRaw.Valid { _ = json.Unmarshal([]byte(prosRaw.String), &z.Pros) }
	if consRaw.Valid { _ = json.Unmarshal([]byte(consRaw.String), &z.Cons) }
	if housingRaw.Valid { _ = json.Unmarshal([]byte(housingRaw.String), &z.HousingTypes) }
	subzones, _ := r.GetSubzonesByZone(ctx, id)
	return &model.DistrictDetail{District: z, Subzones: subzones}, nil
}
