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

-- Бронирования
CREATE TABLE IF NOT EXISTS bookings (
    id              SERIAL PRIMARY KEY,
    apartment_id    INTEGER NOT NULL REFERENCES apartments(id) ON DELETE CASCADE,
    client_id       INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    check_in        DATE NOT NULL,
    check_out       DATE NOT NULL,
    guests_count    INTEGER NOT NULL DEFAULT 1,
    total_price     REAL NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending_approval'
                    CHECK (status IN ('pending_approval','approved','payment_claimed','confirmed','rejected','cancelled')),
    admin_note      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (check_out > check_in)
);
CREATE INDEX IF NOT EXISTS idx_bookings_apartment_id ON bookings(apartment_id);
CREATE INDEX IF NOT EXISTS idx_bookings_client_id ON bookings(client_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
CREATE INDEX IF NOT EXISTS idx_bookings_check_in ON bookings(check_in);

-- Система фильтров
CREATE TABLE IF NOT EXISTS filter_categories (
    id          SERIAL PRIMARY KEY,
    code        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    is_active   INTEGER NOT NULL DEFAULT 1
);
CREATE TABLE IF NOT EXISTS filter_options (
    id          SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES filter_categories(id) ON DELETE CASCADE,
    code        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    is_active   INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX IF NOT EXISTS idx_filter_options_category_id ON filter_options(category_id);
CREATE TABLE IF NOT EXISTS apartment_filters (
    apartment_id     INTEGER NOT NULL REFERENCES apartments(id) ON DELETE CASCADE,
    filter_option_id INTEGER NOT NULL REFERENCES filter_options(id) ON DELETE CASCADE,
    PRIMARY KEY (apartment_id, filter_option_id)
);
CREATE INDEX IF NOT EXISTS idx_apartment_filters_apartment_id ON apartment_filters(apartment_id);
CREATE INDEX IF NOT EXISTS idx_apartment_filters_filter_option_id ON apartment_filters(filter_option_id);

-- Фото квартир
CREATE TABLE IF NOT EXISTS apartment_photos (
    id           SERIAL PRIMARY KEY,
    apartment_id INTEGER NOT NULL REFERENCES apartments(id) ON DELETE CASCADE,
    file_path    TEXT NOT NULL,
    url          TEXT,
    sort_order   INTEGER NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_apartment_photos_apartment_id ON apartment_photos(apartment_id);
CREATE INDEX IF NOT EXISTS idx_apartment_photos_sort_order ON apartment_photos(apartment_id, sort_order);

-- Управление доступностью дат
CREATE TABLE IF NOT EXISTS apartment_availability (
    id             SERIAL PRIMARY KEY,
    apartment_id   INTEGER NOT NULL REFERENCES apartments(id) ON DELETE CASCADE,
    date_from      DATE NOT NULL,
    date_to        DATE NOT NULL,
    status         TEXT NOT NULL DEFAULT 'blocked' CHECK (status IN ('blocked', 'available')),
    note           TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (date_to >= date_from)
);
CREATE INDEX IF NOT EXISTS idx_apt_avail_apartment_id ON apartment_availability(apartment_id);
CREATE INDEX IF NOT EXISTS idx_apt_avail_dates       ON apartment_availability(date_from, date_to);
