CREATE TABLE IF NOT EXISTS zones (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
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

-- Top level districts
INSERT INTO zones (parent_id, city, name, emoji, short_desc, target_audience, pros, cons, housing_types, price_level, best_for, season_note, sort_order)
VALUES
(NULL, 'Одесса', 'Исторический центр', '🏛️',
 'Дерибасовская, Оперный, Потёмкинская лестница',
 'Туристы на 1–3 дня, пары, романтические поездки',
 '[