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

-- Seed apartments (owner = admin user with telegram_id 100000001)
-- Active seed data: only Аркадия (zone_id = 3)

INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, 3,
       'Апартаменты в 5 минутах от пляжа',
       'Современные апартаменты в курортном комплексе. До главного пляжа Одессы — 5 минут пешком. Всё для пляжного отдыха.',
       'ул. Генуэзская, 3',
       1, 'apartments', 2, 2200.00,
       '["photo_ark_1_1.jpg","photo_ark_1_2.jpg"]',
       '["WiFi","Кухня","Кондиционер","Стиральная машина","Сейф"]'
FROM users u WHERE u.telegram_id = 100000001;

INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, 3,
       'Студия с панорамным видом на море',
       'Студия на верхнем этаже с панорамным видом на Чёрное море. Романтика и алые закаты каждый вечер.',
       'Набережная Аркадии, 9',
       1, 'studio', 2, 2600.00,
       '["photo_ark_2_1.jpg","photo_ark_2_2.jpg"]',
       '["WiFi","Кухня","Кондиционер","Балкон","Кофемашина"]'
FROM users u WHERE u.telegram_id = 100000001;

INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, 3,
       'Трёшка для большой компании',
       'Просторная трёхкомнатная квартира для большой компании или семьи. 3 спальни, 2 санузла, большая гостиная с кухней.',
       'ул. Тенистая, 8',
       3, 'family', 6, 4200.00,
       '["photo_ark_3_1.jpg","photo_ark_3_2.jpg","photo_ark_3_3.jpg"]',
       '["WiFi","Кухня","Кондиционер","Стиральная машина","Посудомойка","2 санузла","Детская кроватка","Парковка"]'
FROM users u WHERE u.telegram_id = 100000001
ON CONFLICT DO NOTHING;
