-- ============================================================
-- Seed data - применяется только если SEED=true в .env
-- Все запросы идемпотентны
-- ============================================================

-- Cities
INSERT INTO cities (name, emoji, sort_order) VALUES ('Одесса', '🌊', 1) ON CONFLICT (name) DO NOTHING;
INSERT INTO cities (name, emoji, sort_order) VALUES ('Винница', '🏰', 2) ON CONFLICT (name) DO NOTHING;

-- Zones (Одесса → Аркадия)
INSERT INTO zones (city_id, name, emoji, short_desc, full_desc, target_audience, pros, cons, housing_types, price_level, best_for, season_note, sort_order)
SELECT c.id, 'Аркадия', '🏖', 'Главный пляж, курортная инфраструктура',
    'Курортный район Одессы с лучшими пляжами.',
    'Отдыхающие у моря, пары, семьи',
    '["Пляж рядом","Развитая инфраструктура"]', '["Шумно в пик сезона"]',
    '["Студии","1-2 комнатные","Пентхаусы","Апартаменты"]', 3, 'Пляжный отдых', 'Максимальный спрос летом', 1
FROM cities c WHERE c.name = 'Одесса' ON CONFLICT (city_id, name) DO NOTHING;

-- Subzones (Аркадия)
INSERT INTO subzones (zone_id, name, emoji, short_desc, sort_order)
SELECT z.id, 'Гагарин Плаза', '🏙', 'Деловой и жилой кластер', 1 FROM zones z JOIN cities c ON c.id = z.city_id WHERE z.name = 'Аркадия' AND c.name = 'Одесса'
UNION ALL SELECT z.id, 'Элегия Парк', '🌳', 'Современные ЖК у парка', 2 FROM zones z JOIN cities c ON c.id = z.city_id WHERE z.name = 'Аркадия' AND c.name = 'Одесса'
UNION ALL SELECT z.id, 'Родос / Эллада', '🏢', 'Район башен у моря', 3 FROM zones z JOIN cities c ON c.id = z.city_id WHERE z.name = 'Аркадия' AND c.name = 'Одесса'
UNION ALL SELECT z.id, 'Акрополь', '🏛', 'Премиальная застройка', 4 FROM zones z JOIN cities c ON c.id = z.city_id WHERE z.name = 'Аркадия' AND c.name = 'Одесса'
UNION ALL SELECT z.id, 'Каманина', '🌅', 'Близко к набережной', 5 FROM zones z JOIN cities c ON c.id = z.city_id WHERE z.name = 'Аркадия' AND c.name = 'Одесса'
UNION ALL SELECT z.id, 'Морская сторона', '🌊', 'Линия домов у воды', 6 FROM zones z JOIN cities c ON c.id = z.city_id WHERE z.name = 'Аркадия' AND c.name = 'Одесса'
UNION ALL SELECT z.id, 'Генуэзская', '🚶', 'Центральная улица', 7 FROM zones z JOIN cities c ON c.id = z.city_id WHERE z.name = 'Аркадия' AND c.name = 'Одесса'
UNION ALL SELECT z.id, 'Аркадийская аллея', '🎡', 'Пешеходная зона', 8 FROM zones z JOIN cities c ON c.id = z.city_id WHERE z.name = 'Аркадия' AND c.name = 'Одесса'
UNION ALL SELECT z.id, 'Район Ibiza / Itaka', '🎶', 'Клубная линия', 9 FROM zones z JOIN cities c ON c.id = z.city_id WHERE z.name = 'Аркадия' AND c.name = 'Одесса'
UNION ALL SELECT z.id, 'Тихая Аркадия', '🌙', 'Спокойные улицы', 10 FROM zones z JOIN cities c ON c.id = z.city_id WHERE z.name = 'Аркадия' AND c.name = 'Одесса'
UNION ALL SELECT z.id, 'Ближе к Парку Победы', '🌲', 'У парка', 11 FROM zones z JOIN cities c ON c.id = z.city_id WHERE z.name = 'Аркадия' AND c.name = 'Одесса'
UNION ALL SELECT z.id, 'Ближе к трассе здоровья', '🚴', 'Спортивный маршрут', 12 FROM zones z JOIN cities c ON c.id = z.city_id WHERE z.name = 'Аркадия' AND c.name = 'Одесса'
ON CONFLICT (zone_id, name) DO NOTHING;

-- Zones (Винница)
INSERT INTO zones (city_id, name, emoji, short_desc, sort_order)
SELECT c.id, 'Центр', '🏛', 'Центр города', 1 FROM cities c WHERE c.name = 'Винница' ON CONFLICT (city_id, name) DO NOTHING;
INSERT INTO zones (city_id, name, emoji, short_desc, sort_order)
SELECT c.id, 'Вишенка', '🌳', 'Спальный район', 2 FROM cities c WHERE c.name = 'Винница' ON CONFLICT (city_id, name) DO NOTHING;
INSERT INTO zones (city_id, name, emoji, short_desc, sort_order)
SELECT c.id, 'Замостье', '🌉', 'За рекой', 3 FROM cities c WHERE c.name = 'Винница' ON CONFLICT (city_id, name) DO NOTHING;

-- Subzones (Винница — копии зон для ссылок из apartments)
INSERT INTO subzones (zone_id, name, emoji, short_desc, sort_order)
SELECT z.id, z.name, z.emoji, z.short_desc, z.sort_order FROM zones z JOIN cities c ON c.id = z.city_id WHERE c.name = 'Винница'
ON CONFLICT (zone_id, name) DO NOTHING;

-- Admin user
INSERT INTO users (telegram_id, username, first_name, last_name, role)
VALUES (542389660, 'booking_admin', 'Администратор', 'Системный', 'admin')
ON CONFLICT (telegram_id) DO NOTHING;

-- Client user
INSERT INTO users (telegram_id, username, first_name, last_name, role)
VALUES (7530461559, 'vitaliy', 'Vitaliy', '', 'client')
ON CONFLICT (telegram_id) DO NOTHING;

-- filter_categories
INSERT INTO filter_categories (code, name, sort_order, is_active) VALUES
('location','📍 Локация',1,1),
('sea_distance','🌊 Расстояние до моря',2,1),
('apartment_type','🏠 Тип квартиры',3,1),
('balcony','🏗 Балкон / терраса',4,1),
('sleeping','🛏 Спальные места',5,1),
('electricity','⚡ Электричество и автономность',6,1),
('safety','🔒 Дом и безопасность',7,1)
ON CONFLICT (code) DO NOTHING;

-- filter_options
INSERT INTO filter_options (category_id, code, name, sort_order, is_active) VALUES
((SELECT id FROM filter_categories WHERE code='location'),'zone_gagarin_plaza','Гагарин Плаза',1,1),
((SELECT id FROM filter_categories WHERE code='location'),'zone_elegiya_park','Элегия Парк',2,1),
((SELECT id FROM filter_categories WHERE code='location'),'zone_rodos_ellada','Родос / Эллада',3,1),
((SELECT id FROM filter_categories WHERE code='location'),'zone_akropol','Акрополь',4,1),
((SELECT id FROM filter_categories WHERE code='location'),'zone_kamanina','Каманина',5,1),
((SELECT id FROM filter_categories WHERE code='location'),'zone_morskaya','Морская сторона',6,1),
((SELECT id FROM filter_categories WHERE code='location'),'zone_genuezskaya','Генуэзская',7,1),
((SELECT id FROM filter_categories WHERE code='location'),'zone_arkadiyskaya_alleya','Аркадийская аллея',8,1),
((SELECT id FROM filter_categories WHERE code='location'),'zone_ibiza_itaka','Район Ibiza / Itaka',9,1),
((SELECT id FROM filter_categories WHERE code='location'),'zone_tihaya_arkadiya','Тихая Аркадия',10,1),
((SELECT id FROM filter_categories WHERE code='location'),'zone_park_pobedy','Ближе к Парку Победы',11,1),
((SELECT id FROM filter_categories WHERE code='location'),'zone_trassa_zdorovya','Ближе к трассе здоровья',12,1),
-- Винница
((SELECT id FROM filter_categories WHERE code='location'),'zone_vinnytsia_center','Центр',13,1),
((SELECT id FROM filter_categories WHERE code='location'),'zone_vinnytsia_vishenka','Вишенка',14,1),
((SELECT id FROM filter_categories WHERE code='location'),'zone_vinnytsia_zamostye','Замостье',15,1),
((SELECT id FROM filter_categories WHERE code='sea_distance'),'sea_1_3_min','До моря 1-3 мин',1,1),
((SELECT id FROM filter_categories WHERE code='sea_distance'),'sea_5_min','До моря 5 мин',2,1),
((SELECT id FROM filter_categories WHERE code='sea_distance'),'sea_first_line','Первая линия',3,1),
((SELECT id FROM filter_categories WHERE code='sea_distance'),'sea_direct_view','Вид прямо на море',4,1),
((SELECT id FROM filter_categories WHERE code='sea_distance'),'sea_side_view','Боковой вид на море',5,1),
((SELECT id FROM filter_categories WHERE code='sea_distance'),'view_city','Вид на город',6,1),
((SELECT id FROM filter_categories WHERE code='sea_distance'),'view_sunset','Вид на закат',7,1),
((SELECT id FROM filter_categories WHERE code='sea_distance'),'view_yard','Вид во двор',8,1),
((SELECT id FROM filter_categories WHERE code='sea_distance'),'side_quiet','Тихая сторона',9,1),
((SELECT id FROM filter_categories WHERE code='sea_distance'),'side_south','Южная сторона',10,1),
((SELECT id FROM filter_categories WHERE code='sea_distance'),'side_east','Восточная сторона',11,1),
((SELECT id FROM filter_categories WHERE code='sea_distance'),'side_west','Западная сторона',12,1),
((SELECT id FROM filter_categories WHERE code='apartment_type'),'type_studio','Студия',1,1),
((SELECT id FROM filter_categories WHERE code='apartment_type'),'type_1room','1-комнатная',2,1),
((SELECT id FROM filter_categories WHERE code='apartment_type'),'type_2room','2-комнатная',3,1),
((SELECT id FROM filter_categories WHERE code='apartment_type'),'type_penthouse','Пентхаус',4,1),
((SELECT id FROM filter_categories WHERE code='apartment_type'),'type_apartments','Апартаменты',5,1),
((SELECT id FROM filter_categories WHERE code='apartment_type'),'type_family','Семейная квартира',6,1),
((SELECT id FROM filter_categories WHERE code='balcony'),'has_balcony','Есть балкон',1,1),
((SELECT id FROM filter_categories WHERE code='balcony'),'big_terrace','Большая терраса',2,1),
((SELECT id FROM filter_categories WHERE code='balcony'),'panoramic_windows','Панорамные окна',3,1),
((SELECT id FROM filter_categories WHERE code='balcony'),'smoking_terrace','Можно курить на террасе',4,1),
((SELECT id FROM filter_categories WHERE code='balcony'),'no_smoking','Курение запрещено',5,1),
((SELECT id FROM filter_categories WHERE code='balcony'),'terrace_furniture','Мебель на террасе',6,1),
((SELECT id FROM filter_categories WHERE code='balcony'),'terrace_sunbeds','Лежаки / зона отдыха',7,1),
((SELECT id FROM filter_categories WHERE code='sleeping'),'sleep_double_bed','Двуспальная кровать',1,1),
((SELECT id FROM filter_categories WHERE code='sleeping'),'sleep_king_size','King Size',2,1),
((SELECT id FROM filter_categories WHERE code='sleeping'),'sleep_sofa','Диван',3,1),
((SELECT id FROM filter_categories WHERE code='sleeping'),'sleep_sofa_bed','Раскладной диван',4,1),
((SELECT id FROM filter_categories WHERE code='sleeping'),'sleep_single_beds','Отдельные кровати',5,1),
((SELECT id FROM filter_categories WHERE code='sleeping'),'sleep_child_bed','Детская кровать',6,1),
((SELECT id FROM filter_categories WHERE code='electricity'),'elec_generator','Генератор в доме',1,1),
((SELECT id FROM filter_categories WHERE code='electricity'),'elec_ups','Бесперебойник',2,1),
((SELECT id FROM filter_categories WHERE code='electricity'),'elec_battery','Аккумуляторы',3,1),
((SELECT id FROM filter_categories WHERE code='electricity'),'elec_internet_blackout','Интернет при отключении света',4,1),
((SELECT id FROM filter_categories WHERE code='electricity'),'elec_elevator_blackout','Лифт работает при blackout',5,1),
((SELECT id FROM filter_categories WHERE code='electricity'),'elec_water_blackout','Есть вода при отключении',6,1),
((SELECT id FROM filter_categories WHERE code='safety'),'safety_guard','Охрана',1,1),
((SELECT id FROM filter_categories WHERE code='safety'),'safety_concierge','Консьерж',2,1),
((SELECT id FROM filter_categories WHERE code='safety'),'safety_closed_area','Закрытая территория',3,1),
((SELECT id FROM filter_categories WHERE code='safety'),'safety_cctv','Видеонаблюдение',4,1),
((SELECT id FROM filter_categories WHERE code='safety'),'safety_parking','Паркинг',5,1),
((SELECT id FROM filter_categories WHERE code='safety'),'safety_underground_parking','Подземный паркинг',6,1),
((SELECT id FROM filter_categories WHERE code='safety'),'safety_pets','Можно с животными',7,1),
((SELECT id FROM filter_categories WHERE code='safety'),'safety_self_checkin','Self check-in',8,1)
ON CONFLICT (code) DO NOTHING;

-- Apartments (все 15 тестовых квартир, идемпотентные через WHERE NOT EXISTS)
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM subzones WHERE name='Гагарин Плаза'), 'Студия в Гагарин Плаза с видом на город', 'Светлая студия в новом ЖК Гагарин Плаза. Вид на городские огни, консьерж, охраняемая территория. Генератор в доме — свет есть всегда.', 'просп. Гагарина, 19', 1, 'studio', 2, 1800.00, '[]', '[]' FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Студия в Гагарин Плаза с видом на город');
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM subzones WHERE name='Элегия Парк'), 'Уютная 1-комнатная у парка, тихая сторона', 'Квартира с большим балконом на тихую сторону. Раскладной диван в гостиной, рядом зелёная зона.', 'ул. Парковая, 5', 1, '1room', 3, 1900.00, '[]', '[]' FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Уютная 1-комнатная у парка, тихая сторона');
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM subzones WHERE name='Родос / Эллада'), 'Двушка первая линия — вид прямо на море', 'Квартира в ЖК Родос с прямым видом на Чёрное море. Первая линия, кровать King Size, панорамные окна.', 'Набережная Аркадии, 14а', 2, '2room', 4, 3800.00, '[]', '[]' FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Двушка первая линия — вид прямо на море');
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM subzones WHERE name='Акрополь'), 'Пентхаус с террасой — вид на закат', 'Роскошный пентхаус с большой террасой. Вид на закат, бесперебойник, видеонаблюдение. Self check-in.', 'ул. Акропольская, 3', 2, 'penthouse', 4, 6500.00, '[]', '[]' FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Пентхаус с террасой — вид на закат');
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM subzones WHERE name='Каманина'), 'Студия 2 мин до пляжа, боковой вид на море', 'Компактная студия в 2 минутах ходьбы до пляжа. Балкон с боковым видом на море.', 'ул. Каманина, 16', 1, 'studio', 2, 2100.00, '[]', '[]' FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Студия 2 мин до пляжа, боковой вид на море');
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM subzones WHERE name='Морская сторона'), 'Апартаменты первая линия — панорама моря', 'Апартаменты прямо у воды с панорамным видом. Большая терраса с лежаками. Генератор, аккумуляторы.', 'Набережная Аркадии, 22', 2, 'apartments', 4, 4500.00, '[]', '[]' FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Апартаменты первая линия — панорама моря');
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM subzones WHERE name='Генуэзская'), '1-комнатная в закрытом ЖК на Генуэзской', 'Квартира в закрытом ЖК. Охрана, паркинг, консьерж. Две отдельные кровати.', 'ул. Генуэзская, 24б', 1, '1room', 2, 2000.00, '[]', '[]' FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = '1-комнатная в закрытом ЖК на Генуэзской');
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM subzones WHERE name='Аркадийская аллея'), 'Семейная двушка у аллеи — для семьи с детьми', 'Просторная 2-комнатная рядом с Аркадийской аллеей. Детская кроватка, можно с животными.', 'Аркадийская аллея, 7', 2, '2room', 4, 2800.00, '[]', '[]' FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Семейная двушка у аллеи — для семьи с детьми');
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM subzones WHERE name='Район Ibiza / Itaka'), 'Студия у Ibiza — 5 мин до пляжа', 'Стильная студия рядом с клубной зоной. До пляжа 5 минут, южная сторона. Self check-in.', 'ул. Генуэзская, 30', 1, 'studio', 2, 2300.00, '[]', '[]' FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Студия у Ibiza — 5 мин до пляжа');
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM subzones WHERE name='Тихая Аркадия'), 'Тихая Аркадия — семейная квартира', 'Спокойный уголок Аркадии вдали от клубов. Большой диван, лифт работает при отключении света.', 'ул. Тихая, 12', 2, 'family', 5, 2400.00, '[]', '[]' FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Тихая Аркадия — семейная квартира');
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM subzones WHERE name='Ближе к Парку Победы'), 'Двушка у Парка Победы — автономная', 'Квартира рядом с Парком Победы. Аккумуляторы, вода и интернет при отключении света. Паркинг.', 'просп. Шевченко, 4', 2, '2room', 4, 2200.00, '[]', '[]' FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Двушка у Парка Победы — автономная');
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM subzones WHERE name='Ближе к трассе здоровья'), 'Апартаменты у трассы здоровья — для активных', 'Современная квартира для активного отдыха. Панорамные окна, кровать King Size. Подземный паркинг.', 'ул. Черноморская, 8', 1, '1room', 2, 2100.00, '[]', '[]' FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Апартаменты у трассы здоровья — для активных');

-- Три уникальные квартиры (subzone_id=3 — Аркадия)
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, 3, 'Апартаменты в 5 минутах от пляжа', 'Современные апартаменты в курортном комплексе. До главного пляжа Одессы — 5 минут пешком.', 'ул. Генуэзская, 3', 1, 'apartments', 2, 2200.00, '["photo_ark_1_1.jpg","photo_ark_1_2.jpg"]', '["WiFi","Кухня","Кондиционер","Стиральная машина","Сейф"]'
FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Апартаменты в 5 минутах от пляжа');
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, 3, 'Студия с панорамным видом на море', 'Студия на верхнем этаже с панорамным видом на Чёрное море. Романтика и алые закаты каждый вечер.', 'Набережная Аркадии, 9', 1, 'studio', 2, 2600.00, '["photo_ark_2_1.jpg","photo_ark_2_2.jpg"]', '["WiFi","Кухня","Кондиционер","Балкон","Кофемашина"]'
FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Студия с панорамным видом на море');
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, 3, 'Трёшка для большой компании', 'Просторная трёхкомнатная квартира для большой компании или семьи. 3 спальни, 2 санузла.', 'ул. Тенистая, 8', 3, 'family', 6, 4200.00, '["photo_ark_3_1.jpg","photo_ark_3_2.jpg","photo_ark_3_3.jpg"]', '["WiFi","Кухня","Кондиционер","Стиральная машина","Посудомойка","2 санузла","Детская кроватка","Парковка"]'
FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Трёшка для большой компании');

-- Винница: тестовая квартира
INSERT INTO apartments (owner_id, subzone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT s.id FROM subzones s JOIN zones z ON z.id = s.zone_id JOIN cities c ON c.id = z.city_id WHERE s.name='Центр' AND c.name='Винница'), 'Студия в центре Винницы', 'Уютная студия в самом центре города. Рядом парк, фонтаны и кафе.', 'ул. Соборная, 15', 1, 'studio', 2, 1200.00, '[]', '["WiFi","Кондиционер","Кухня"]'
FROM users u WHERE u.telegram_id = 542389660 AND NOT EXISTS (SELECT 1 FROM apartments WHERE title = 'Студия в центре Винницы');

-- Вычисляем city_id для всех квартир (если ещё не заполнен)
UPDATE apartments a SET city_id = c.id
FROM subzones s JOIN zones z ON z.id = s.zone_id JOIN cities c ON c.id = z.city_id
WHERE a.subzone_id = s.id AND a.city_id IS NULL;

-- apartment_filters для ВСЕХ тестовых квартир
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Студия в Гагарин Плаза с видом на город' AND fo.code IN ('zone_gagarin_plaza','view_city','elec_generator','safety_guard','safety_concierge') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Уютная 1-комнатная у парка, тихая сторона' AND fo.code IN ('zone_elegiya_park','side_quiet','has_balcony','sleep_sofa_bed','safety_pets') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Двушка первая линия — вид прямо на море' AND fo.code IN ('zone_rodos_ellada','sea_first_line','sea_direct_view','sleep_king_size','panoramic_windows','safety_underground_parking') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Пентхаус с террасой — вид на закат' AND fo.code IN ('zone_akropol','view_sunset','big_terrace','terrace_furniture','elec_ups','safety_cctv','safety_self_checkin') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Студия 2 мин до пляжа, боковой вид на море' AND fo.code IN ('zone_kamanina','sea_1_3_min','sea_side_view','has_balcony','smoking_terrace') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Апартаменты первая линия — панорама моря' AND fo.code IN ('zone_morskaya','sea_first_line','sea_direct_view','big_terrace','terrace_sunbeds','elec_generator','elec_battery','elec_internet_blackout') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='1-комнатная в закрытом ЖК на Генуэзской' AND fo.code IN ('zone_genuezskaya','safety_closed_area','safety_parking','safety_guard','safety_concierge','sleep_single_beds') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Семейная двушка у аллеи — для семьи с детьми' AND fo.code IN ('zone_arkadiyskaya_alleya','view_yard','sleep_double_bed','sleep_child_bed','safety_pets') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Студия у Ibiza — 5 мин до пляжа' AND fo.code IN ('zone_ibiza_itaka','sea_5_min','side_south','safety_self_checkin','no_smoking') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Тихая Аркадия — семейная квартира' AND fo.code IN ('zone_tihaya_arkadiya','side_east','sleep_sofa','sleep_child_bed','elec_elevator_blackout') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Двушка у Парка Победы — автономная' AND fo.code IN ('zone_park_pobedy','side_west','elec_battery','elec_water_blackout','elec_internet_blackout','safety_parking') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Апартаменты у трассы здоровья — для активных' AND fo.code IN ('zone_trassa_zdorovya','panoramic_windows','sleep_king_size','safety_underground_parking') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Апартаменты в 5 минутах от пляжа' AND fo.code IN ('type_apartments') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Студия с панорамным видом на море' AND fo.code IN ('type_studio','has_balcony') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Трёшка для большой компании' AND fo.code IN ('type_family','sleep_child_bed') ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.title='Студия в центре Винницы' AND fo.code IN ('zone_vinnytsia_center','type_studio') ON CONFLICT DO NOTHING;

-- type_* для всех квартир (по apartment_type)
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.apartment_type = 'studio' AND fo.code = 'type_studio' ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.apartment_type = '1room' AND fo.code = 'type_1room' ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.apartment_type = '2room' AND fo.code = 'type_2room' ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.apartment_type = 'penthouse' AND fo.code = 'type_penthouse' ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.apartment_type = 'apartments' AND fo.code = 'type_apartments' ON CONFLICT DO NOTHING;
INSERT INTO apartment_filters (apartment_id, filter_option_id) SELECT a.id, fo.id FROM apartments a, filter_options fo WHERE a.apartment_type = 'family' AND fo.code = 'type_family' ON CONFLICT DO NOTHING;

-- ============================================================
-- Тестовые данные: занятость квартир (apartment_availability)
-- ============================================================

-- Квартира id=4 (Студия Гагарин Плаза): заблокированные диапазоны (ремонт, личные нужды)
INSERT INTO apartment_availability (apartment_id, date_from, date_to, status, note)
SELECT a.id,
       (CURRENT_DATE + i * 14)::DATE,
       (CURRENT_DATE + i * 14 + 4)::DATE,
       'blocked',
       CASE i WHEN 0 THEN 'Ремонт' WHEN 1 THEN 'Личные нужды' ELSE 'Недоступно' END
FROM apartments a, (VALUES (0),(1),(2)) t(i)
WHERE a.title = 'Студия в Гагарин Плаза с видом на город'
  AND NOT EXISTS (
    SELECT 1 FROM apartment_availability aa
    WHERE aa.apartment_id = a.id AND aa.status = 'blocked'
  );

-- Квартира «Двушка первая линия»: заблокирован целый диапазон (бронь хозяина)
INSERT INTO apartment_availability (apartment_id, date_from, date_to, status, note)
SELECT a.id,
       (CURRENT_DATE + 7)::DATE,
       (CURRENT_DATE + 12)::DATE,
       'blocked',
       'Хозяйская бронь'
FROM apartments a
WHERE a.title = 'Двушка первая линия — вид прямо на море'
  AND NOT EXISTS (
    SELECT 1 FROM apartment_availability aa
    WHERE aa.apartment_id = a.id AND aa.status = 'blocked'
  );

-- Квартира «Пентхаус»: явно открытый период (available) + заблокированный
INSERT INTO apartment_availability (apartment_id, date_from, date_to, status, note)
SELECT a.id, d_from::DATE, d_to::DATE, st, note
FROM apartments a,
     (VALUES
       ((CURRENT_DATE + 3)::TEXT,  (CURRENT_DATE + 6)::TEXT,  'blocked',   'Технический перерыв'),
       ((CURRENT_DATE + 20)::TEXT, (CURRENT_DATE + 25)::TEXT, 'blocked',   'Ремонт балкона'),
       ((CURRENT_DATE + 7)::TEXT,  (CURRENT_DATE + 19)::TEXT, 'available', 'Открыто для бронирования')
     ) t(d_from, d_to, st, note)
WHERE a.title = 'Пентхаус с террасой — вид на закат'
  AND NOT EXISTS (
    SELECT 1 FROM apartment_availability aa WHERE aa.apartment_id = a.id
  );

-- Тестовые бронирования (bookings) — имитация реальных заявок
INSERT INTO bookings (apartment_id, client_id, check_in, check_out, guests_count, total_price, status)
SELECT a.id, u.id,
       (CURRENT_DATE + 2)::DATE,
       (CURRENT_DATE + 5)::DATE,
       2, a.price_per_night * 3, 'confirmed'
FROM apartments a, users u
WHERE a.title = 'Студия в Гагарин Плаза с видом на город'
  AND u.telegram_id = 7530461559
  AND NOT EXISTS (
    SELECT 1 FROM bookings b WHERE b.apartment_id = a.id AND b.status = 'confirmed'
  );

INSERT INTO bookings (apartment_id, client_id, check_in, check_out, guests_count, total_price, status)
SELECT a.id, u.id,
       (CURRENT_DATE + 15)::DATE,
       (CURRENT_DATE + 19)::DATE,
       3, a.price_per_night * 4, 'approved'
FROM apartments a, users u
WHERE a.title = 'Двушка первая линия — вид прямо на море'
  AND u.telegram_id = 7530461559
  AND NOT EXISTS (
    SELECT 1 FROM bookings b WHERE b.apartment_id = a.id AND b.status = 'approved'
  );









