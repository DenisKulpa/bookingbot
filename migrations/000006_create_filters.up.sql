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
