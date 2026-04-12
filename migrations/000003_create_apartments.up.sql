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
    is_available    BOOLEAN NOT NULL DEFAULT true,
    created_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at      DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_apartments_zone_id ON apartments(zone_id);
CREATE INDEX IF NOT EXISTS idx_apartments_owner_id ON apartments(owner_id);
CREATE INDEX IF NOT EXISTS idx_apartments_is_available ON apartments(is_available);
CREATE INDEX IF NOT EXISTS idx_apartments_price ON apartments(price_per_night);

-- Seed apartments (owner = admin user with telegram_id 100000001)
-- Zone IDs mirror 000001_create_zones: 1=Исторический центр, 2=Приморский,
--   3=Аркадия, 4=Фонтанка, 5=Молдаванка, 6=Черёмушки/Таирова

-- Исторический центр (zone_id = 1) — 3 квартиры
INSERT OR IGNORE INTO apartments (owner_id, zone_id, title, description, address, rooms, max_guests, price_per_night, photos, amenities)
SELECT u.id, 1,
       'Апартаменты на Дерибасовской',
       'Уютные апартаменты в самом сердце Одессы, два шага от Дерибасовской. Кирпичный дом начала XX века, высокие потолки, дубовый паркет.',
       'ул. Дерибасовская, 14',
       1, 2, 2500.00,
       '["photo_center_1_1.jpg","photo_center_1_2.jpg"]',
       '["WiFi","Кухня","Кондиционер","Стиральная машина"]'
FROM users u WHERE u.telegram_id = 100000001;

INSERT OR IGNORE INTO apartments (owner_id, zone_id, title, description, address, rooms, max_guests, price_per_night, photos, amenities)
SELECT u.id, 1,
       'Студия у Оперного театра',
       'Светлая студия с видом на исторические улочки, 3 минуты пешком до Оперного театра. Свежий ремонт, современная мебель.',
       'ул. Пушкинская, 7',
       1, 2, 2200.00,
       '["photo_center_2_1.jpg","photo_center_2_2.jpg"]',
       '["WiFi","Кухня","Кондиционер","Рабочее место"]'
FROM users u WHERE u.telegram_id = 100000001;

INSERT OR IGNORE INTO apartments (owner_id, zone_id, title, description, address, rooms, max_guests, price_per_night, photos, amenities)
SELECT u.id, 1,
       'Двухкомнатная на Приморском бульваре',
       'Просторная квартира с потрясающим видом на море с Приморского бульвара. Антикварная мебель, атмосфера старой Одессы.',
       'Приморский бульвар, 5',
       2, 4, 3800.00,
       '["photo_center_3_1.jpg","photo_center_3_2.jpg","photo_center_3_3.jpg"]',
       '["WiFi","Кухня","Кондиционер","Балкон с видом на море","Стиральная машина","Посудомойка"]'
FROM users u WHERE u.telegram_id = 100000001;

-- Приморский (zone_id = 2) — 2 квартиры
INSERT OR IGNORE INTO apartments (owner_id, zone_id, title, description, address, rooms, max_guests, price_per_night, photos, amenities)
SELECT u.id, 2,
       'Уютная студия у парка Шевченко',
       'Тихая студия в зелёном Приморском районе, 5 минут до парка Шевченко и пляжа Ланжерон.',
       'ул. Парковая, 12',
       1, 2, 1800.00,
       '["photo_prim_1_1.jpg","photo_prim_1_2.jpg"]',
       '["WiFi","Кухня","Кондиционер","Гладильная доска"]'
FROM users u WHERE u.telegram_id = 100000001;

INSERT OR IGNORE INTO apartments (owner_id, zone_id, title, description, address, rooms, max_guests, price_per_night, photos, amenities)
SELECT u.id, 2,
       'Светлая двушка на Французском бульваре',
       'Современная двухкомнатная квартира в новом ЖК на Французском бульваре. Закрытый двор, подземный паркинг, охрана.',
       'Французский бульвар, 44',
       2, 4, 2800.00,
       '["photo_prim_2_1.jpg","photo_prim_2_2.jpg","photo_prim_2_3.jpg"]',
       '["WiFi","Кухня","Кондиционер","Балкон","Стиральная машина","Парковка"]'
FROM users u WHERE u.telegram_id = 100000001;

-- Аркадия (zone_id = 3) — 3 квартиры
INSERT OR IGNORE INTO apartments (owner_id, zone_id, title, description, address, rooms, max_guests, price_per_night, photos, amenities)
SELECT u.id, 3,
       'Апартаменты в 5 минутах от пляжа',
       'Современные апартаменты в курортном комплексе. До главного пляжа Одессы — 5 минут пешком. Всё для пляжного отдыха.',
       'ул. Генуэзская, 3',
       1, 2, 2200.00,
       '["photo_ark_1_1.jpg","photo_ark_1_2.jpg"]',
       '["WiFi","Кухня","Кондиционер","Стиральная машина","Сейф"]'
FROM users u WHERE u.telegram_id = 100000001;

INSERT OR IGNORE INTO apartments (owner_id, zone_id, title, description, address, rooms, max_guests, price_per_night, photos, amenities)
SELECT u.id, 3,
       'Студия с панорамным видом на море',
       'Студия на верхнем этаже с панорамным видом на Чёрное море. Романтика и алые закаты каждый вечер.',
       'Набережная Аркадии, 9',
       1, 2, 2600.00,
       '["photo_ark_2_1.jpg","photo_ark_2_2.jpg"]',
       '["WiFi","Кухня","Кондиционер","Балкон","Кофемашина"]'
FROM users u WHERE u.telegram_id = 100000001;

INSERT OR IGNORE INTO apartments (owner_id, zone_id, title, description, address, rooms, max_guests, price_per_night, photos, amenities)
SELECT u.id, 3,
       'Трёшка для большой компании',
       'Просторная трёхкомнатная квартира для большой компании или семьи. 3 спальни, 2 санузла, большая гостиная с кухней.',
       'ул. Тенистая, 8',
       3, 6, 4200.00,
       '["photo_ark_3_1.jpg","photo_ark_3_2.jpg","photo_ark_3_3.jpg"]',
       '["WiFi","Кухня","Кондиционер","Стиральная машина","Посудомойка","2 санузла","Детская кроватка","Парковка"]'
FROM users u WHERE u.telegram_id = 100000001;

-- Фонтанка (zone_id = 4) — 2 квартиры
INSERT OR IGNORE INTO apartments (owner_id, zone_id, title, description, address, rooms, max_guests, price_per_night, photos, amenities)
SELECT u.id, 4,
       'Дача у моря — Малый Фонтан',
       'Уютный дачный домик с летней верандой и небольшим садом. До пляжа — 10 минут пешком. Мангал, тишина, свежий воздух.',
       'пер. Дачный, 5, Малый Фонтан',
       2, 4, 1200.00,
       '["photo_font_1_1.jpg","photo_font_1_2.jpg"]',
       '["WiFi","Кухня","Мангал","Летняя веранда","Парковка во дворе"]'
FROM users u WHERE u.telegram_id = 100000001;

INSERT OR IGNORE INTO apartments (owner_id, zone_id, title, description, address, rooms, max_guests, price_per_night, photos, amenities)
SELECT u.id, 4,
       'Тихий домик на Большом Фонтане',
       'Отдельный домик в тишине Большого Фонтана. Дикий пляж в 7 минутах ходьбы. Идеально для отдыха от городской суеты.',
       'ул. Черноморская, 21, Большой Фонтан',
       2, 4, 1000.00,
       '["photo_font_2_1.jpg","photo_font_2_2.jpg"]',
       '["WiFi","Кухня","Мангал","Терраса","Велосипеды"]'
FROM users u WHERE u.telegram_id = 100000001;

-- Молдаванка (zone_id = 5) — 2 квартиры
INSERT OR IGNORE INTO apartments (owner_id, zone_id, title, description, address, rooms, max_guests, price_per_night, photos, amenities)
SELECT u.id, 5,
       'Колоритная квартира в одесском дворике',
       'Настоящая Одесса: знаменитые дворики-колодцы, балкончик с геранью, соседи с характером. Дом начала XX века.',
       'ул. Мясоедовская, 12',
       1, 2, 950.00,
       '["photo_mold_1_1.jpg","photo_mold_1_2.jpg"]',
       '["WiFi","Кухня","Стиральная машина","Балкон"]'
FROM users u WHERE u.telegram_id = 100000001;

INSERT OR IGNORE INTO apartments (owner_id, zone_id, title, description, address, rooms, max_guests, price_per_night, photos, amenities)
SELECT u.id, 5,
       'Бюджетная студия рядом с Привозом',
       'Компактная и чистая студия в 10 минутах пешком от рынка Привоз. Всё необходимое для комфортного проживания.',
       'ул. Торговая, 3',
       1, 2, 800.00,
       '["photo_mold_2_1.jpg","photo_mold_2_2.jpg"]',
       '["WiFi","Кухня","Кондиционер"]'
FROM users u WHERE u.telegram_id = 100000001;

-- Черёмушки / Таирова (zone_id = 6) — 2 квартиры
INSERT OR IGNORE INTO apartments (owner_id, zone_id, title, description, address, rooms, max_guests, price_per_night, photos, amenities)
SELECT u.id, 6,
       'Современная двушка в Таирово',
       'Квартира в современном ЖК района Таирово. Свежий ремонт, вся бытовая техника, детская площадка во дворе, охраняемая парковка.',
       'пр. Академика Глушко, 7',
       2, 4, 1100.00,
       '["photo_tair_1_1.jpg","photo_tair_1_2.jpg","photo_tair_1_3.jpg"]',
       '["WiFi","Кухня","Кондиционер","Стиральная машина","Посудомойка","Парковка"]'
FROM users u WHERE u.telegram_id = 100000001;

INSERT OR IGNORE INTO apartments (owner_id, zone_id, title, description, address, rooms, max_guests, price_per_night, photos, amenities)
SELECT u.id, 6,
       'Просторная трёшка в Черёмушках',
       'Большая трёхкомнатная квартира для семьи с детьми. Спокойный спальный район, вся инфраструктура рядом, отличная транспортная доступность.',
       'ул. Марсельская, 32',
       3, 5, 1400.00,
       '["photo_cher_1_1.jpg","photo_cher_1_2.jpg"]',
       '["WiFi","Кухня","Кондиционер","Стиральная машина","Детская кроватка","Парковка"]'
FROM users u WHERE u.telegram_id = 100000001;
