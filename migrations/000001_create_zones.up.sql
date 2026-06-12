-- Города
CREATE TABLE IF NOT EXISTS cities (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    emoji      TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active  INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Районы (внутри города)
CREATE TABLE IF NOT EXISTS zones (
    id              SERIAL PRIMARY KEY,
    city_id         INTEGER NOT NULL REFERENCES cities(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    emoji           TEXT,
    short_desc      TEXT,
    full_desc       TEXT,
    target_audience TEXT,
    pros            TEXT,
    cons            TEXT,
    housing_types   TEXT,
    price_level     INTEGER CHECK (price_level BETWEEN 1 AND 3),
    best_for        TEXT,
    season_note     TEXT,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    is_active       INTEGER NOT NULL DEFAULT 1,
    UNIQUE (city_id, name)
);
CREATE INDEX IF NOT EXISTS idx_zones_city_id ON zones(city_id);

-- Подрайоны / микрорайоны
CREATE TABLE IF NOT EXISTS subzones (
    id         SERIAL PRIMARY KEY,
    zone_id    INTEGER NOT NULL REFERENCES zones(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    emoji      TEXT,
    short_desc TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active  INTEGER NOT NULL DEFAULT 1,
    UNIQUE (zone_id, name)
);
CREATE INDEX IF NOT EXISTS idx_subzones_zone_id ON subzones(zone_id);
