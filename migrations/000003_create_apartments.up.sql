CREATE TABLE IF NOT EXISTS apartments (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id        INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    zone_id         INTEGER REFERENCES zones(id) ON DELETE SET NULL,
    title           TEXT NOT NULL,
    description     TEXT,
    address         TEXT,
    rooms           INTEGER NOT NULL DEFAULT 1,
    max_guests      INTEGER NOT NULL DEFAULT 2,
    price_per_night REAL NOT NULL,
    photos          TEXT NOT NULL DEFAULT '[]',  -- JSON array of file paths
    amenities       TEXT NOT NULL DEFAULT '[]',  -- JSON array of strings
    is_available    INTEGER NOT NULL DEFAULT 1,
    created_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at      DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_apartments_zone_id ON apartments(zone_id);
CREATE INDEX IF NOT EXISTS idx_apartments_owner_id ON apartments(owner_id);
CREATE INDEX IF NOT EXISTS idx_apartments_is_available ON apartments(is_available);
CREATE INDEX IF NOT EXISTS idx_apartments_price ON apartments(price_per_night);
