CREATE TABLE IF NOT EXISTS bookings (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    apartment_id    INTEGER NOT NULL REFERENCES apartments(id) ON DELETE CASCADE,
    client_id       INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    check_in        DATE NOT NULL,
    check_out       DATE NOT NULL,
    guests_count    INTEGER NOT NULL DEFAULT 1,
    total_price     REAL NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending_approval'
                    CHECK (status IN (
                        'pending_approval',
                        'approved',
                        'payment_claimed',
                        'confirmed',
                        'rejected',
                        'cancelled'
                    )),
    admin_note      TEXT,
    created_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    CHECK (check_out > check_in)
);

CREATE INDEX IF NOT EXISTS idx_bookings_apartment_id ON bookings(apartment_id);
CREATE INDEX IF NOT EXISTS idx_bookings_client_id ON bookings(client_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
CREATE INDEX IF NOT EXISTS idx_bookings_check_in ON bookings(check_in);
