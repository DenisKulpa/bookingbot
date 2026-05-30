-- ============================================================
-- Seed: тестовые квартиры по всем подзонам Аркадии с фильтрами
-- Подзоны (parent_id=3, subzones inserted in 000001):
--   Гагарин Плаза, Элегия Парк, Родос/Эллада, Акрополь,
--   Каманина, Морская сторона, Генуэзская, Аркадийская аллея,
--   Район Ibiza/Itaka, Тихая Аркадия, Парк Победы, Трасса здоровья
-- ============================================================

-- 1. Гагарин Плаза — студия, вид на город, генератор, охрана
INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM zones WHERE name = 'Гагарин Плаза'),
    'Студия в Гагарин Плаза с видом на город',
    'Светлая студия в новом ЖК Гагарин Плаза. Вид на городские огни, консьерж, охраняемая территория. Генератор в доме — свет есть всегда.',
    'просп. Гагарина, 19',
    1, 'studio', 2, 1800.00,
    '[]', '["WiFi","Кондиционер","Кухня","Консьерж","Охрана","Генератор"]'
FROM users u WHERE u.telegram_id = 100000001;

-- 2. Элегия Парк — 1-комнатная, тихая сторона, балкон, раскладной диван
INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM zones WHERE name = 'Элегия Парк'),
    'Уютная 1-комнатная у парка, тихая сторона',
    'Квартира с большим балконом на тихую сторону. Раскладной диван в гостиной, рядом зелёная зона. Тихий двор, можно с животными.',
    'ул. Парковая, 5',
    1, '1room', 3, 1900.00,
    '[]', '["WiFi","Кондиционер","Кухня","Балкон","Стиральная машина","Можно с животными"]'
FROM users u WHERE u.telegram_id = 100000001;

-- 3. Родос / Эллада — 2-комнатная, вид прямо на море, первая линия, King Size
INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM zones WHERE name = 'Родос / Эллада'),
    'Двушка первая линия — вид прямо на море',
    'Квартира в ЖК Родос с прямым видом на Чёрное море. Первая линия от воды, кровать King Size, панорамные окна. Подземный паркинг.',
    'Набережная Аркадии, 14а',
    2, '2room', 4, 3800.00,
    '[]', '["WiFi","Кондиционер","Кухня","Панорамные окна","King Size","Подземный паркинг","Посудомойка"]'
FROM users u WHERE u.telegram_id = 100000001;

-- 4. Акрополь — пентхаус, терраса, закат, бесперебойник, видеонаблюдение
INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM zones WHERE name = 'Акрополь'),
    'Пентхаус с террасой — вид на закат',
    'Роскошный пентхаус с большой террасой с мебелью. Вид на закат над морем, бесперебойник, видеонаблюдение в доме. Self check-in.',
    'ул. Акропольская, 3',
    2, 'penthouse', 4, 6500.00,
    '[]', '["WiFi","Кондиционер","Кухня","Большая терраса","Мебель на террасе","Вид на закат","Бесперебойник","Видеонаблюдение","Self check-in"]'
FROM users u WHERE u.telegram_id = 100000001;

-- 5. Каманина — студия, до моря 1-3 мин, боковой вид на море, курить на террасе
INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM zones WHERE name = 'Каманина'),
    'Студия 2 мин до пляжа, боковой вид на море',
    'Компактная студия в 2 минутах ходьбы до пляжа. Балкон с боковым видом на море, можно курить на террасе. Всё для пляжного отдыха.',
    'ул. Каманина, 16',
    1, 'studio', 2, 2100.00,
    '[]', '["WiFi","Кондиционер","Кухня","Балкон","Кофемашина"]'
FROM users u WHERE u.telegram_id = 100000001;

-- 6. Морская сторона — апартаменты, первая линия, прямой вид на море, лежаки на террасе
INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM zones WHERE name = 'Морская сторона'),
    'Апартаменты первая линия — панорама моря',
    'Апартаменты прямо у воды с панорамным видом на море. Большая терраса с лежаками. Генератор, аккумуляторы, интернет при отключении света.',
    'Набережная Аркадии, 22',
    2, 'apartments', 4, 4500.00,
    '[]', '["WiFi","Кондиционер","Кухня","Большая терраса","Лежаки / зона отдыха","Генератор","Аккумуляторы","Интернет при отключении света"]'
FROM users u WHERE u.telegram_id = 100000001;

-- 7. Генуэзская — 1-комнатная, закрытая территория, паркинг, отдельные кровати
INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM zones WHERE name = 'Генуэзская'),
    '1-комнатная в закрытом ЖК на Генуэзской',
    'Квартира в закрытом жилом комплексе. Охрана, паркинг, консьерж. Две отдельные кровати — идеально для коллег или друзей.',
    'ул. Генуэзская, 24б',
    1, '1room', 2, 2000.00,
    '[]', '["WiFi","Кондиционер","Кухня","Закрытая территория","Паркинг","Охрана","Консьерж","Отдельные кровати"]'
FROM users u WHERE u.telegram_id = 100000001;

-- 8. Аркадийская аллея — 2-комнатная, вид во двор, двуспальная + детская кровать, можно с животными
INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM zones WHERE name = 'Аркадийская аллея'),
    'Семейная двушка у аллеи — для семьи с детьми',
    'Просторная двухкомнатная квартира рядом с Аркадийской аллеей. Детская кроватка, можно с животными. Тихий двор. Стиральная машина.',
    'Аркадийская аллея, 7',
    2, '2room', 4, 2800.00,
    '[]', '["WiFi","Кондиционер","Кухня","Двуспальная кровать","Детская кровать","Стиральная машина","Можно с животными"]'
FROM users u WHERE u.telegram_id = 100000001;

-- 9. Район Ibiza / Itaka — студия, до моря 5 мин, южная сторона, self check-in, курение запрещено
INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM zones WHERE name = 'Район Ibiza / Itaka'),
    'Студия у Ibiza — 5 мин до пляжа',
    'Стильная студия рядом с клубной зоной. До пляжа 5 минут, солнечная южная сторона. Self check-in, курение в квартире запрещено.',
    'ул. Генуэзская, 30',
    1, 'studio', 2, 2300.00,
    '[]', '["WiFi","Кондиционер","Кухня","Self check-in","Курение запрещено","Кофемашина"]'
FROM users u WHERE u.telegram_id = 100000001;

-- 10. Тихая Аркадия — семейная квартира, диван, восточная сторона, лифт при blackout
INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM zones WHERE name = 'Тихая Аркадия'),
    'Тихая Аркадия — семейная квартира',
    'Спокойный уголок Аркадии вдали от клубов. Большой диван, лифт работает при отключении света. Восточная сторона — утреннее солнце.',
    'ул. Тихая, 12',
    2, 'family', 5, 2400.00,
    '[]', '["WiFi","Кондиционер","Кухня","Диван","Стиральная машина","Детская кровать","Лифт работает при blackout"]'
FROM users u WHERE u.telegram_id = 100000001;

-- 11. Ближе к Парку Победы — 2-комнатная, западная сторона, вода при отключении, паркинг
INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM zones WHERE name = 'Ближе к Парку Победы'),
    'Двушка у Парка Победы — автономная',
    'Квартира рядом с Парком Победы. Полная автономность: аккумуляторы, вода при отключении, интернет при blackout. Паркинг.',
    'просп. Шевченко, 4',
    2, '2room', 4, 2200.00,
    '[]', '["WiFi","Кондиционер","Кухня","Аккумуляторы","Есть вода при отключении","Интернет при отключении света","Паркинг","Стиральная машина"]'
FROM users u WHERE u.telegram_id = 100000001;

-- 12. Ближе к трассе здоровья — 1-комнатная, King Size, панорамные окна, подземный паркинг
INSERT INTO apartments (owner_id, zone_id, title, description, address, rooms, apartment_type, max_guests, price_per_night, photos, amenities)
SELECT u.id, (SELECT id FROM zones WHERE name = 'Ближе к трассе здоровья'),
    'Апартаменты у трассы здоровья — для активных',
    'Современная квартира для любителей активного отдыха. Панорамные окна, кровать King Size. Подземный паркинг, велосипедная стоянка рядом.',
    'ул. Черноморская, 8',
    1, '1room', 2, 2100.00,
    '[]', '["WiFi","Кондиционер","Кухня","Панорамные окна","King Size","Подземный паркинг","Стиральная машина"]'
FROM users u WHERE u.telegram_id = 100000001;


-- ============================================================
-- Привязка квартир к filter_options через apartment_filters
-- ============================================================

-- helper: получаем id квартир по названию и связываем с фильтрами

-- 1. Студия Гагарин Плаза → location:zone_gagarin_plaza, view:city, elec:generator, safety:guard, safety:concierge
INSERT INTO apartment_filters (apartment_id, filter_option_id)
SELECT a.id, fo.id FROM apartments a, filter_options fo
WHERE a.title = 'Студия в Гагарин Плаза с видом на город'
  AND fo.code IN ('zone_gagarin_plaza','view_city','elec_generator','safety_guard','safety_concierge')
ON CONFLICT DO NOTHING;

-- 2. Элегия Парк → zone_elegiya_park, side_quiet, has_balcony, sleep_sofa_bed, safety_pets
INSERT INTO apartment_filters (apartment_id, filter_option_id)
SELECT a.id, fo.id FROM apartments a, filter_options fo
WHERE a.title = 'Уютная 1-комнатная у парка, тихая сторона'
  AND fo.code IN ('zone_elegiya_park','side_quiet','has_balcony','sleep_sofa_bed','safety_pets')
ON CONFLICT DO NOTHING;

-- 3. Родос / Эллада → zone_rodos_ellada, sea_first_line, sea_direct_view, sleep_king_size, panoramic_windows, safety_underground_parking
INSERT INTO apartment_filters (apartment_id, filter_option_id)
SELECT a.id, fo.id FROM apartments a, filter_options fo
WHERE a.title = 'Двушка первая линия — вид прямо на море'
  AND fo.code IN ('zone_rodos_ellada','sea_first_line','sea_direct_view','sleep_king_size','panoramic_windows','safety_underground_parking')
ON CONFLICT DO NOTHING;

-- 4. Акрополь → zone_akropol, view_sunset, big_terrace, terrace_furniture, elec_ups, safety_cctv, safety_self_checkin
INSERT INTO apartment_filters (apartment_id, filter_option_id)
SELECT a.id, fo.id FROM apartments a, filter_options fo
WHERE a.title = 'Пентхаус с террасой — вид на закат'
  AND fo.code IN ('zone_akropol','view_sunset','big_terrace','terrace_furniture','elec_ups','safety_cctv','safety_self_checkin')
ON CONFLICT DO NOTHING;

-- 5. Каманина → zone_kamanina, sea_1_3_min, sea_side_view, has_balcony, smoking_terrace
INSERT INTO apartment_filters (apartment_id, filter_option_id)
SELECT a.id, fo.id FROM apartments a, filter_options fo
WHERE a.title = 'Студия 2 мин до пляжа, боковой вид на море'
  AND fo.code IN ('zone_kamanina','sea_1_3_min','sea_side_view','has_balcony','smoking_terrace')
ON CONFLICT DO NOTHING;

-- 6. Морская сторона → zone_morskaya, sea_first_line, sea_direct_view, big_terrace, terrace_sunbeds, elec_generator, elec_battery, elec_internet_blackout
INSERT INTO apartment_filters (apartment_id, filter_option_id)
SELECT a.id, fo.id FROM apartments a, filter_options fo
WHERE a.title = 'Апартаменты первая линия — панорама моря'
  AND fo.code IN ('zone_morskaya','sea_first_line','sea_direct_view','big_terrace','terrace_sunbeds','elec_generator','elec_battery','elec_internet_blackout')
ON CONFLICT DO NOTHING;

-- 7. Генуэзская → zone_genuezskaya, safety_closed_area, safety_parking, safety_guard, safety_concierge, sleep_single_beds
INSERT INTO apartment_filters (apartment_id, filter_option_id)
SELECT a.id, fo.id FROM apartments a, filter_options fo
WHERE a.title = '1-комнатная в закрытом ЖК на Генуэзской'
  AND fo.code IN ('zone_genuezskaya','safety_closed_area','safety_parking','safety_guard','safety_concierge','sleep_single_beds')
ON CONFLICT DO NOTHING;

-- 8. Аркадийская аллея → zone_arkadiyskaya_alleya, view_yard, sleep_double_bed, sleep_child_bed, safety_pets
INSERT INTO apartment_filters (apartment_id, filter_option_id)
SELECT a.id, fo.id FROM apartments a, filter_options fo
WHERE a.title = 'Семейная двушка у аллеи — для семьи с детьми'
  AND fo.code IN ('zone_arkadiyskaya_alleya','view_yard','sleep_double_bed','sleep_child_bed','safety_pets')
ON CONFLICT DO NOTHING;

-- 9. Ibiza/Itaka → zone_ibiza_itaka, sea_5_min, side_south, safety_self_checkin, no_smoking
INSERT INTO apartment_filters (apartment_id, filter_option_id)
SELECT a.id, fo.id FROM apartments a, filter_options fo
WHERE a.title = 'Студия у Ibiza — 5 мин до пляжа'
  AND fo.code IN ('zone_ibiza_itaka','sea_5_min','side_south','safety_self_checkin','no_smoking')
ON CONFLICT DO NOTHING;

-- 10. Тихая Аркадия → zone_tihaya_arkadiya, side_east, sleep_sofa, sleep_child_bed, elec_elevator_blackout
INSERT INTO apartment_filters (apartment_id, filter_option_id)
SELECT a.id, fo.id FROM apartments a, filter_options fo
WHERE a.title = 'Тихая Аркадия — семейная квартира'
  AND fo.code IN ('zone_tihaya_arkadiya','side_east','sleep_sofa','sleep_child_bed','elec_elevator_blackout')
ON CONFLICT DO NOTHING;

-- 11. Парк Победы → zone_park_pobedy, side_west, elec_battery, elec_water_blackout, elec_internet_blackout, safety_parking
INSERT INTO apartment_filters (apartment_id, filter_option_id)
SELECT a.id, fo.id FROM apartments a, filter_options fo
WHERE a.title = 'Двушка у Парка Победы — автономная'
  AND fo.code IN ('zone_park_pobedy','side_west','elec_battery','elec_water_blackout','elec_internet_blackout','safety_parking')
ON CONFLICT DO NOTHING;

-- 12. Трасса здоровья → zone_trassa_zdorovya, panoramic_windows, sleep_king_size, safety_underground_parking
INSERT INTO apartment_filters (apartment_id, filter_option_id)
SELECT a.id, fo.id FROM apartments a, filter_options fo
WHERE a.title = 'Апартаменты у трассы здоровья — для активных'
  AND fo.code IN ('zone_trassa_zdorovya','panoramic_windows','sleep_king_size','safety_underground_parking')
ON CONFLICT DO NOTHING;
