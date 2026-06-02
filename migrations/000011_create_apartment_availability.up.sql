-- Таблица управления доступностью квартир.
-- Арендодатель может вручную закрывать/открывать диапазоны дат.
-- Приоритет: если дата есть в bookings (approved/confirmed) — занята независимо от этой таблицы.

CREATE TABLE IF NOT EXISTS apartment_availability (
    id             SERIAL PRIMARY KEY,
    apartment_id   INTEGER NOT NULL REFERENCES apartments(id) ON DELETE CASCADE,
    date_from      DATE NOT NULL,
    date_to        DATE NOT NULL,   -- включительно
    status         TEXT NOT NULL DEFAULT 'blocked'
                   CHECK (status IN ('blocked', 'available')),
    note           TEXT,            -- комментарий арендодателя (напр. "Ремонт", "Личные нужды")
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (date_to >= date_from)
);

CREATE INDEX IF NOT EXISTS idx_apt_avail_apartment_id ON apartment_availability(apartment_id);
CREATE INDEX IF NOT EXISTS idx_apt_avail_dates       ON apartment_availability(date_from, date_to);
