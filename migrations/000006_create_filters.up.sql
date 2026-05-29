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

INSERT INTO filter_categories (code, name, sort_order, is_active)
VALUES
('location', 'Локация в Аркадии', 1, 1),
('sea_distance', 'Расстояние до моря', 2, 1),
('apartment_type', 'Тип квартиры', 3, 1),
('balcony', 'Балкон / терраса', 4, 1),
('sleeping', 'Спальные места', 5, 1),
('electricity', 'Электричество и автономность', 6, 1),
('safety', 'Дом и безопасность', 7, 1)
ON CONFLICT (code) DO NOTHING;

INSERT INTO filter_options (category_id, code, name, sort_order, is_active)
VALUES
((SELECT id FROM filter_categories WHERE code = 'location'), 'zone_gagarin_plaza', 'Гагарин Плаза', 1, 1),
((SELECT id FROM filter_categories WHERE code = 'location'), 'zone_elegiya_park', 'Элегия Парк', 2, 1),
((SELECT id FROM filter_categories WHERE code = 'location'), 'zone_rodos_ellada', 'Родос / Эллада', 3, 1),
((SELECT id FROM filter_categories WHERE code = 'location'), 'zone_akropol', 'Акрополь', 4, 1),
((SELECT id FROM filter_categories WHERE code = 'location'), 'zone_kamanina', 'Каманина', 5, 1),
((SELECT id FROM filter_categories WHERE code = 'location'), 'zone_morskaya', 'Морская сторона', 6, 1),
((SELECT id FROM filter_categories WHERE code = 'location'), 'zone_genuezskaya', 'Генуэзская', 7, 1),
((SELECT id FROM filter_categories WHERE code = 'location'), 'zone_arkadiyskaya_alleya', 'Аркадийская аллея', 8, 1),
((SELECT id FROM filter_categories WHERE code = 'location'), 'zone_ibiza_itaka', 'Район Ibiza / Itaka', 9, 1),
((SELECT id FROM filter_categories WHERE code = 'location'), 'zone_tihaya_arkadiya', 'Тихая Аркадия', 10, 1),
((SELECT id FROM filter_categories WHERE code = 'location'), 'zone_park_pobedy', 'Ближе к Парку Победы', 11, 1),
((SELECT id FROM filter_categories WHERE code = 'location'), 'zone_trassa_zdorovya', 'Ближе к трассе здоровья', 12, 1),

((SELECT id FROM filter_categories WHERE code = 'sea_distance'), 'sea_1_3_min', 'До моря 1–3 мин', 1, 1),
((SELECT id FROM filter_categories WHERE code = 'sea_distance'), 'sea_5_min', 'До моря 5 мин', 2, 1),
((SELECT id FROM filter_categories WHERE code = 'sea_distance'), 'sea_first_line', 'Первая линия', 3, 1),
((SELECT id FROM filter_categories WHERE code = 'sea_distance'), 'sea_direct_view', 'Вид прямо на море', 4, 1),
((SELECT id FROM filter_categories WHERE code = 'sea_distance'), 'sea_side_view', 'Боковой вид на море', 5, 1),
((SELECT id FROM filter_categories WHERE code = 'sea_distance'), 'view_city', 'Вид на город', 6, 1),
((SELECT id FROM filter_categories WHERE code = 'sea_distance'), 'view_sunset', 'Вид на закат', 7, 1),
((SELECT id FROM filter_categories WHERE code = 'sea_distance'), 'view_yard', 'Вид во двор', 8, 1),
((SELECT id FROM filter_categories WHERE code = 'sea_distance'), 'side_quiet', 'Тихая сторона', 9, 1),
((SELECT id FROM filter_categories WHERE code = 'sea_distance'), 'side_south', 'Южная сторона', 10, 1),
((SELECT id FROM filter_categories WHERE code = 'sea_distance'), 'side_east', 'Восточная сторона', 11, 1),
((SELECT id FROM filter_categories WHERE code = 'sea_distance'), 'side_west', 'Западная сторона', 12, 1),

((SELECT id FROM filter_categories WHERE code = 'apartment_type'), 'type_studio', 'Студия', 1, 1),
((SELECT id FROM filter_categories WHERE code = 'apartment_type'), 'type_1room', '1-комнатная', 2, 1),
((SELECT id FROM filter_categories WHERE code = 'apartment_type'), 'type_2room', '2-комнатная', 3, 1),
((SELECT id FROM filter_categories WHERE code = 'apartment_type'), 'type_penthouse', 'Пентхаус', 4, 1),
((SELECT id FROM filter_categories WHERE code = 'apartment_type'), 'type_apartments', 'Апартаменты', 5, 1),
((SELECT id FROM filter_categories WHERE code = 'apartment_type'), 'type_family', 'Семейная квартира', 6, 1),

((SELECT id FROM filter_categories WHERE code = 'balcony'), 'has_balcony', 'Есть балкон', 1, 1),
((SELECT id FROM filter_categories WHERE code = 'balcony'), 'big_terrace', 'Большая терраса', 2, 1),
((SELECT id FROM filter_categories WHERE code = 'balcony'), 'panoramic_windows', 'Панорамные окна', 3, 1),
((SELECT id FROM filter_categories WHERE code = 'balcony'), 'smoking_terrace', 'Можно курить на террасе', 4, 1),
((SELECT id FROM filter_categories WHERE code = 'balcony'), 'no_smoking', 'Курение запрещено', 5, 1),
((SELECT id FROM filter_categories WHERE code = 'balcony'), 'terrace_furniture', 'Мебель на террасе', 6, 1),
((SELECT id FROM filter_categories WHERE code = 'balcony'), 'terrace_sunbeds', 'Лежаки / зона отдыха', 7, 1),

((SELECT id FROM filter_categories WHERE code = 'sleeping'), 'sleep_double_bed', 'Двуспальная кровать', 1, 1),
((SELECT id FROM filter_categories WHERE code = 'sleeping'), 'sleep_king_size', 'King Size', 2, 1),
((SELECT id FROM filter_categories WHERE code = 'sleeping'), 'sleep_sofa', 'Диван', 3, 1),
((SELECT id FROM filter_categories WHERE code = 'sleeping'), 'sleep_sofa_bed', 'Раскладной диван', 4, 1),
((SELECT id FROM filter_categories WHERE code = 'sleeping'), 'sleep_single_beds', 'Отдельные кровати', 5, 1),
((SELECT id FROM filter_categories WHERE code = 'sleeping'), 'sleep_child_bed', 'Детская кровать', 6, 1),

((SELECT id FROM filter_categories WHERE code = 'electricity'), 'elec_generator', 'Генератор в доме', 1, 1),
((SELECT id FROM filter_categories WHERE code = 'electricity'), 'elec_ups', 'Бесперебойник', 2, 1),
((SELECT id FROM filter_categories WHERE code = 'electricity'), 'elec_battery', 'Аккумуляторы', 3, 1),
((SELECT id FROM filter_categories WHERE code = 'electricity'), 'elec_internet_blackout', 'Интернет при отключении света', 4, 1),
((SELECT id FROM filter_categories WHERE code = 'electricity'), 'elec_elevator_blackout', 'Лифт работает при blackout', 5, 1),
((SELECT id FROM filter_categories WHERE code = 'electricity'), 'elec_water_blackout', 'Есть вода при отключении', 6, 1),

((SELECT id FROM filter_categories WHERE code = 'safety'), 'safety_guard', 'Охрана', 1, 1),
((SELECT id FROM filter_categories WHERE code = 'safety'), 'safety_concierge', 'Консьерж', 2, 1),
((SELECT id FROM filter_categories WHERE code = 'safety'), 'safety_closed_area', 'Закрытая территория', 3, 1),
((SELECT id FROM filter_categories WHERE code = 'safety'), 'safety_cctv', 'Видеонаблюдение', 4, 1),
((SELECT id FROM filter_categories WHERE code = 'safety'), 'safety_parking', 'Паркинг', 5, 1),
((SELECT id FROM filter_categories WHERE code = 'safety'), 'safety_underground_parking', 'Подземный паркинг', 6, 1),
((SELECT id FROM filter_categories WHERE code = 'safety'), 'safety_pets', 'Можно с животными', 7, 1),
((SELECT id FROM filter_categories WHERE code = 'safety'), 'safety_self_checkin', 'Self check-in', 8, 1)
ON CONFLICT (code) DO NOTHING;
