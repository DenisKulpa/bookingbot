CREATE TABLE IF NOT EXISTS apartment_photos (
    id           SERIAL PRIMARY KEY,
    apartment_id INTEGER NOT NULL REFERENCES apartments(id) ON DELETE CASCADE,
    file_path    TEXT NOT NULL,          -- относительный путь: uploads/apartments/{id}/filename.jpg
    url          TEXT,                   -- публичный URL для отдачи через HTTP
    sort_order   INTEGER NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_apartment_photos_apartment_id ON apartment_photos(apartment_id);
CREATE INDEX IF NOT EXISTS idx_apartment_photos_sort_order ON apartment_photos(apartment_id, sort_order);
