CREATE TABLE IF NOT EXISTS zones (
    id              SERIAL PRIMARY KEY,
    parent_id       INTEGER REFERENCES zones(id) ON DELETE CASCADE,
    city            TEXT NOT NULL DEFAULT 'Одесса',
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
    is_active       INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_zones_parent_id ON zones(parent_id);
CREATE INDEX IF NOT EXISTS idx_zones_is_active ON zones(is_active);

-- Защита от дубликатов: одна зона — одно имя
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'uq_zones_name' AND conrelid = 'zones'::regclass
    ) THEN
        ALTER TABLE zones ADD CONSTRAINT uq_zones_name UNIQUE (name);
    END IF;
END $$;
