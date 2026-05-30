CREATE TABLE IF NOT EXISTS apartments (
    id               SERIAL PRIMARY KEY,
    owner_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    zone_id          INTEGER REFERENCES zones(id) ON DELETE SET NULL,
    title            TEXT NOT NULL,
    description      TEXT,
    address          TEXT,
    rooms            INTEGER NOT NULL DEFAULT 1,
    apartment_type   TEXT CHECK (apartment_type IN ('studio','1room','2room','penthouse','apartments','family')) DEFAULT NULL,
    max_guests       INTEGER NOT NULL DEFAULT 2,
    price_per_night  REAL NOT NULL,
    photos           TEXT NOT NULL DEFAULT '[]',  -- JSON array of file paths
    amenities        TEXT NOT NULL DEFAULT '[]',  -- JSON array of strings
    is_available     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_apartments_zone_id ON apartments(zone_id);
CREATE INDEX IF NOT EXISTS idx_apartments_owner_id ON apartments(owner_id);
CREATE INDEX IF NOT EXISTS idx_apartments_is_available ON apartments(is_available);
CREATE INDEX IF NOT EXISTS idx_apartments_price ON apartments(price_per_night);
