package telegram

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/DenisKulpa/bookingbot/internal/repository"
)

const (
	callbackDistrict       = "district:"
	callbackApartment      = "apartment:"
	callbackBack           = "back:districts"
	callbackFilterCat      = "filter_cat:"
	callbackFilter         = "filter:"
	callbackBackFilters    = "back:filters"
	callbackToggleFilter   = "toggle:"
	callbackReset          = "reset:filters"
	callbackSearch         = "search:filters"
	callbackBooking        = "booking:"
	callbackBookingConfirm = "booking_confirm:"
	callbackBookingCancel  = "booking_cancel:"
	callbackBookingApprove = "booking_approve:"
	callbackBookingReject  = "booking_reject:"
	callbackAskQuestion    = "ask_question:"
	callbackChatReply      = "chat_reply:"
)

// ─── Filter data ─────────────────────────────────────────────────────────[...]

type filterOption struct {
	Code  string
	Label string
}

type filterCategory struct {
	Code    string
	Label   string
	Options []filterOption
}

var filterCategories = []filterCategory{
	{Code: "location", Label: "📍 Локация", Options: []filterOption{
		{"zone_gagarin_plaza", "Гагарин Плаза"},
		{"zone_elegiya_park", "Элегия Парк"},
		{"zone_rodos_ellada", "Родос / Эллада"},
		{"zone_akropol", "Акрополь"},
		{"zone_kamanina", "Каманина"},
		{"zone_morskaya", "Морская сторона"},
		{"zone_genuezskaya", "Генуэзская"},
		{"zone_arkadiyskaya_alleya", "Аркадийская аллея"},
		{"zone_ibiza_itaka", "Район Ibiza / Itaka"},
		{"zone_tihaya_arkadiya", "Тихая Аркадия"},
		{"zone_park_pobedy", "Ближе к Парку Победы"},
		{"zone_trassa_zdorovya", "Ближе к трассе здоровья"},
	}},
	{Code: "sea_distance", Label: "🌊 Расстояние до моря", Options: []filterOption{
		{"sea_1_3_min", "До моря 1–3 мин"},
		{"sea_5_min", "До моря 5 мин"},
		{"sea_first_line", "Первая линия"},
		{"sea_direct_view", "Вид прямо на море"},
		{"sea_side_view", "Боковой вид на море"},
		{"view_city", "Вид на город"},
		{"view_sunset", "Вид на закат"},
		{"view_yard", "Вид во двор"},
		{"side_quiet", "Тихая сторона"},
		{"side_south", "Южная сторона"},
		{"side_east", "Восточная сторона"},
		{"side_west", "Западная сторона"},
	}},
	{Code: "apartment_type", Label: "🏠 Тип квартиры", Options: []filterOption{
		{"type_studio", "Студия"},
		{"type_1room", "1-комнатная"},
		{"type_2room", "2-комнатная"},
		{"type_penthouse", "Пентхаус"},
		{"type_apartments", "Апартаменты"},
		{"type_family", "Семейная квартира"},
	}},
	{Code: "balcony", Label: "🏗 Балкон / терраса", Options: []filterOption{
		{"has_balcony", "Есть балкон"},
		{"big_terrace", "Большая терраса"},
		{"panoramic_windows", "Панорамные окна"},
		{"smoking_terrace", "Можно курить на террасе"},
		{"no_smoking", "Курение запрещено"},
		{"terrace_furniture", "Мебель на террасе"},
		{"terrace_sunbeds", "Лежаки / зона отдыха"},
	}},
	{Code: "sleeping", Label: "🛏 Спальные места", Options: []filterOption{
		{"sleep_double_bed", "Двуспальная кровать"},
		{"sleep_king_size", "King Size"},
		{"sleep_sofa", "Диван"},
		{"sleep_sofa_bed", "Раскладной диван"},
		{"sleep_single_beds", "Отдельные кровати"},
		{"sleep_child_bed", "Детская кровать"},
	}},
	{Code: "electricity", Label: "⚡ Электричество и автономность", Options: []filterOption{
		{"elec_generator", "Генератор в доме"},
		{"elec_ups", "Бесперебойник"},
		{"elec_battery", "Аккумуляторы"},
		{"elec_internet_blackout", "Интернет при отключении света"},
		{"elec_elevator_blackout", "Лифт работает при blackout"},
		{"elec_water_blackout", "Есть вода при отключении"},
	}},
	{Code: "safety", Label: "🔒 Дом и безопасность", Options: []filterOption{
		{"safety_guard", "Охрана"},
		{"safety_concierge", "Консьерж"},
		{"safety_closed_area", "Закрытая территория"},
		{"safety_cctv", "Видеонаблюдение"},
		{"safety_parking", "Паркинг"},
		{"safety_underground_parking", "Подземный паркинг"},
		{"safety_pets", "Можно с животными"},
		{"safety_self_checkin", "Self check-in"},
	}},
}

func filterListText(filters map[string]bool) string {
	count := 0
	for _, v := range filters {
		if v {
			count++
		}
	}
	if count == 0 {
		return "🏖 *Поиск жилья в Аркадии по фильтрам*\n\nВыберите категорию:"
	}
	return fmt.Sprintf("🏖 *Поиск жилья в Аркадии по фильтрам*\n\nВыбрано фильтров: *%d*\nВыберите категорию:", count)
}

func filterListRows(filters map[string]bool) [][]tgbotapi.InlineKeyboardButton {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, cat := range filterCategories {
		hasActive := false
		for _, opt := range cat.Options {
			if filters[opt.Code] {
				hasActive = true
				break
			}
		}
		label := cat.Label
		if hasActive {
			label = label + " ✅"
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, callbackFilterCat+cat.Code),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔍 Показать результаты", callbackSearch),
	))
	hasAny := false
	for _, v := range filters {
		if v {
			hasAny = true
			break
		}
	}
	if hasAny {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✖️ Сбросить фильтры", callbackReset),
		))
	}
	return rows
}

// ─── Bot ──────────────────────────────────────────────────────────────────────

type Bot struct {
	client      *Client
	zoneRepo    *repository.ZoneRepository
	aptRepo     *repository.ApartmentRepository
	photoRepo   *repository.PhotoRepository
	bookingRepo *repository.BookingRepository
	userRepo    *repository.UserRepository
	uploadsRoot string
	mu          sync.Mutex
	sessions    map[int64]*session
}

func NewBot(client *Client, zoneRepo *repository.ZoneRepository, aptRepo *repository.ApartmentRepository, photoRepo *repository.PhotoRepository, bookingRepo *repository.BookingRepository, userRepo *repository.UserRepository, uploadsRoot string) *Bot {
	return &Bot{
		client:      client,
		zoneRepo:    zoneRepo,
		aptRepo:     aptRepo,
		photoRepo:   photoRepo,
		bookingRepo: bookingRepo,
		userRepo:    userRepo,
		uploadsRoot: uploadsRoot,
		sessions:    make(map[int64]*session),
	}
}

type session struct {
	filters map[string]bool
	// Бронирование
	bookingStep  int // 0=нет, 1=выбор check-in, 2=выбор check-out
	bookingAptID int
	bookingIn    time.Time
	bookingOut   time.Time
	// Чат
	chatAptID       int   // > 0 — клиент в режиме чата с владельцем этой квартиры
	chatOwnerTgID   int64 // telegram_id владельца квартиры
	chatAptTitle    string
	replyToClientID int64 // > 0 — арендодатель отвечает этому клиенту (chat_id)
}

func (b *Bot) getSession(chatID int64) *session {
	b.mu.Lock()
	defer b.mu.Unlock()
	if s, ok := b.sessions[chatID]; ok {
		return s
	}
	s := &session{filters: make(map[string]bool)}
	b.sessions[chatID] = s
	return s
}

func (b *Bot) toggleFilter(chatID int64, code string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := b.sessions[chatID]
	if s == nil {
		s = &session{filters: make(map[string]bool)}
		b.sessions[chatID] = s
	}
	s.filters[code] = !s.filters[code]
}

func (b *Bot) resetFilters(chatID int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.sessions[chatID] = &session{filters: make(map[string]bool)}
}

func (b *Bot) activeFilters(chatID int64) map[string]bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := b.sessions[chatID]
	if s == nil {
		return map[string]bool{}
	}
	result := make(map[string]bool, len(s.filters))
	for k, v := range s.filters {
		result[k] = v
	}
	return result
}

// Run starts polling for updates and blocking until context is done.
func (b *Bot) Run(ctx context.Context) {
	updates, err := b.client.GetUpdates(60)
	if err != nil {
		log.Printf("bot: failed to get updates channel: %v", err)
		return
	}

	log.Println("bot: polling started")

	for {
		select {
		case <-ctx.Done():
			log.Println("bot: polling stopped")
			return
		case update, ok := <-updates:
			if !ok {
				return
			}
			b.handleUpdate(ctx, &update)
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, u *tgbotapi.Update) {
	if u.Message != nil {
		b.handleMessage(ctx, u.Message)
	} else if u.CallbackQuery != nil {
		b.handleCallback(ctx, u.CallbackQuery)
	}
}

// ─── Message commands ────────────────────────────────────────────────────────

func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	if !msg.IsCommand() {
		b.handleChatMessage(ctx, msg)
		return
	}

	switch msg.Command() {
	case "start":
		b.cmdStart(ctx, msg)
	case "search":
		b.cmdSearch(ctx, msg)
	case "help":
		b.cmdHelp(ctx, msg)
	case "stopchat":
		b.cmdStopChat(msg)
	default:
		_ = b.client.SendMessage(msg.Chat.ID, "Неизвестная команда. Используйте /help для списка команд.")
	}
}

func (b *Bot) cmdStart(ctx context.Context, msg *tgbotapi.Message) {
	name := msg.From.FirstName
	if name == "" {
		name = msg.From.UserName
	}
	text := fmt.Sprintf("👋 Привет, *%s*\\!\n\nЯ помогу найти и забронировать квартиру в Одессе\\.", escapeMarkdownV2(name))
	_ = b.client.sendMarkdownV2(msg.Chat.ID, text)
	b.sendFilterList(msg.Chat.ID)
}

func (b *Bot) cmdHelp(_ context.Context, msg *tgbotapi.Message) {
	text := "*Список команд:*\n\n" +
		"/start — приветствие\n" +
		"/search — поиск жилья в Аркадии по фильтрам\n" +
		"/help — это сообщение"
	_ = b.client.sendMarkdownV2(msg.Chat.ID, text)
}

func (b *Bot) cmdSearch(_ context.Context, msg *tgbotapi.Message) {
	b.sendFilterList(msg.Chat.ID)
}

// ─── Callbacks ───────────────────────────────────────────────────────────────

func (b *Bot) handleCallback(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	_ = b.client.AnswerCallbackQuery(cq.ID, "", false)

	data := cq.Data
	chatID := cq.Message.Chat.ID
	msgID := cq.Message.MessageID

	switch {
	// ── New filter flow ──
	case data == callbackBackFilters:
		b.editFilterList(chatID, msgID)

	case data == callbackReset:
		b.resetFilters(chatID)
		b.editFilterList(chatID, msgID)

	case data == callbackSearch:
		b.editSearchResults(ctx, chatID, msgID)

	case strings.HasPrefix(data, callbackToggleFilter):
		code := strings.TrimPrefix(data, callbackToggleFilter)
		b.toggleFilter(chatID, code)
		catCode := b.findCategoryByFilterCode(code)
		b.editFilterOptions(chatID, msgID, catCode)

	case strings.HasPrefix(data, callbackFilterCat):
		catCode := strings.TrimPrefix(data, callbackFilterCat)
		b.editFilterOptions(chatID, msgID, catCode)

	// ── Legacy district flow (kept for future use) ──
	case data == callbackBack:
		b.editDistrictList(ctx, chatID, msgID)

	case strings.HasPrefix(data, callbackDistrict):
		idStr := strings.TrimPrefix(data, callbackDistrict)
		var id int
		fmt.Sscanf(idStr, "%d", &id)
		b.editDistrictDetail(ctx, chatID, msgID, id)

	case strings.HasPrefix(data, callbackApartment):
		idStr := strings.TrimPrefix(data, callbackApartment)
		var id int
		fmt.Sscanf(idStr, "%d", &id)
		b.editApartmentDetail(ctx, chatID, msgID, id)

	// ── Booking flow ──
	case data == callbackCalIgnore:
		// ничего не делаем — кнопка-заглушка

	case strings.HasPrefix(data, callbackBooking):
		idStr := strings.TrimPrefix(data, callbackBooking)
		aptID, _ := strconv.Atoi(idStr)
		b.startBookingCheckIn(chatID, msgID, aptID)

	case strings.HasPrefix(data, callbackCalNav):
		// cal_nav:YYYY-MM:APTID
		b.handleCalNav(chatID, msgID, cq.From, strings.TrimPrefix(data, callbackCalNav))

	case strings.HasPrefix(data, callbackCalDay):
		// cal_day:YYYY-MM-DD:APTID
		b.handleCalDay(ctx, chatID, msgID, cq.From, strings.TrimPrefix(data, callbackCalDay))

	case strings.HasPrefix(data, callbackBookingConfirm):
		b.confirmBooking(ctx, chatID, msgID, cq.From, strings.TrimPrefix(data, callbackBookingConfirm))

	case strings.HasPrefix(data, callbackBookingCancel):
		aptID, _ := strconv.Atoi(strings.TrimPrefix(data, callbackBookingCancel))
		b.cancelBookingFlow(chatID, aptID)
		b.editApartmentDetail(ctx, chatID, msgID, aptID)

	case strings.HasPrefix(data, callbackBookingApprove):
		b.handleBookingApprove(ctx, chatID, msgID, strings.TrimPrefix(data, callbackBookingApprove))

	case strings.HasPrefix(data, callbackBookingReject):
		b.handleBookingReject(ctx, chatID, msgID, strings.TrimPrefix(data, callbackBookingReject))

	case strings.HasPrefix(data, callbackAskQuestion):
		aptID, _ := strconv.Atoi(strings.TrimPrefix(data, callbackAskQuestion))
		b.handleAskQuestion(ctx, chatID, aptID)

	case strings.HasPrefix(data, callbackChatReply):
		clientChatID, _ := strconv.ParseInt(strings.TrimPrefix(data, callbackChatReply), 10, 64)
		b.handleChatReplyMode(chatID, clientChatID)
	}
}

// ─── Filter screens ───────────────────────────────────────────────────────────

func (b *Bot) sendFilterList(chatID int64) {
	filters := b.activeFilters(chatID)
	_ = b.client.SendMessageWithKeyboard(chatID, filterListText(filters), filterListRows(filters))
}

func (b *Bot) editFilterList(chatID int64, msgID int) {
	filters := b.activeFilters(chatID)
	_ = b.client.EditMessage(chatID, msgID, filterListText(filters), filterListRows(filters))
}

func (b *Bot) editFilterOptions(chatID int64, msgID int, catCode string) {
	var cat *filterCategory
	for i := range filterCategories {
		if filterCategories[i].Code == catCode {
			cat = &filterCategories[i]
			break
		}
	}
	if cat == nil {
		return
	}
	filters := b.activeFilters(chatID)
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, opt := range cat.Options {
		var label string
		if filters[opt.Code] {
			label = "✔️  " + opt.Label
		} else {
			label = "➖  " + opt.Label
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, callbackToggleFilter+opt.Code),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к фильтрам", callbackBackFilters),
	))
	_ = b.client.EditMessage(chatID, msgID, cat.Label+"\n\nВыберите подходящие опции:", rows)
}

func (b *Bot) findCategoryByFilterCode(filterCode string) string {
	for _, cat := range filterCategories {
		for _, opt := range cat.Options {
			if opt.Code == filterCode {
				return cat.Code
			}
		}
	}
	return ""
}

func (b *Bot) editSearchResults(ctx context.Context, chatID int64, msgID int) {
	filters := b.activeFilters(chatID)

	var selected []string
	var selectedCodes []string
	for _, cat := range filterCategories {
		for _, opt := range cat.Options {
			if filters[opt.Code] {
				selected = append(selected, opt.Label)
				selectedCodes = append(selectedCodes, opt.Code)
			}
		}
	}

	apts, err := b.aptRepo.GetByFilters(ctx, selectedCodes)
	if err != nil || len(apts) == 0 {
		_ = b.client.EditMessage(chatID, msgID,
			"😔 По выбранным фильтрам квартир не найдено.",
			[][]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к фильтрам", callbackBackFilters),
				),
			},
		)
		return
	}

	var sb strings.Builder
	if len(selected) > 0 {
		sb.WriteString(fmt.Sprintf("🔍 *Результаты поиска*\n\nФильтры: %s\n\n", strings.Join(selected, " • ")))
	} else {
		sb.WriteString("🔍 *Все доступные квартиры в Аркадии:*\n\n")
	}
	sb.WriteString(fmt.Sprintf("Найдено квартир: *%d*", len(apts)))

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, apt := range apts {
		label := fmt.Sprintf("🏠 %s — %.0f грн/ночь", apt.Title, apt.PricePerNight)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, callbackApartment+fmt.Sprint(apt.ID)),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к фильтрам", callbackBackFilters),
	))
	_ = b.client.EditMessage(chatID, msgID, sb.String(), rows)
}

// ─── Legacy district screens (kept for future use) ───────────────────────────

func (b *Bot) sendDistrictList(ctx context.Context, chatID int64, _ int) {
	zones, err := b.zoneRepo.GetTopLevel(ctx)
	if err != nil {
		log.Printf("bot: GetTopLevel: %v", err)
		_ = b.client.SendMessage(chatID, "Ошибка загрузки районов.")
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, z := range zones {
		label := z.Name
		if z.Emoji != "" {
			label = z.Emoji + " " + z.Name
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, callbackDistrict+fmt.Sprint(z.ID)),
		))
	}

	_ = b.client.SendMessageWithKeyboard(chatID, "🏙 *Выберите район Одессы:*", rows)
}

func (b *Bot) editDistrictList(ctx context.Context, chatID int64, msgID int) {
	zones, err := b.zoneRepo.GetTopLevel(ctx)
	if err != nil {
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, z := range zones {
		label := z.Name
		if z.Emoji != "" {
			label = z.Emoji + " " + z.Name
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, callbackDistrict+fmt.Sprint(z.ID)),
		))
	}

	_ = b.client.EditMessage(chatID, msgID, "🏙 *Выберите район Одессы:*", rows)
}

func (b *Bot) editDistrictDetail(ctx context.Context, chatID int64, msgID int, districtID int) {
	detail, err := b.zoneRepo.GetDistrictDetail(ctx, districtID)
	if err != nil {
		log.Printf("bot: GetDistrictDetail(%d): %v", districtID, err)
		return
	}
	d := detail.District

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%s %s*\n", d.Emoji, d.Name))
	if d.ShortDesc != "" {
		sb.WriteString(fmt.Sprintf("_%s_\n", d.ShortDesc))
	}
	if d.TargetAudience != "" {
		sb.WriteString(fmt.Sprintf("\n👤 *Для кого:* %s\n", d.TargetAudience))
	}
	if len(d.Pros) > 0 {
		sb.WriteString("\n✅ *Плюсы:*\n")
		for _, p := range d.Pros {
			sb.WriteString("• " + p + "\n")
		}
	}
	if len(d.Cons) > 0 {
		sb.WriteString("\n❌ *Минусы:*\n")
		for _, c := range d.Cons {
			sb.WriteString("• " + c + "\n")
		}
	}
	if d.BestFor != "" {
		sb.WriteString(fmt.Sprintf("\n🎯 *Лучше всего для:* %s\n", d.BestFor))
	}
	priceStar := strings.Repeat("💰", d.PriceLevel)
	if priceStar != "" {
		sb.WriteString(fmt.Sprintf("\n%s уровень цен\n", priceStar))
	}

	apts, _ := b.aptRepo.GetByZone(ctx, districtID, true)
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, apt := range apts {
		label := fmt.Sprintf("🏠 %s — %.0f грн/ночь", apt.Title, apt.PricePerNight)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, callbackApartment+fmt.Sprint(apt.ID)),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к районам", callbackBack),
	))

	_ = b.client.EditMessage(chatID, msgID, sb.String(), rows)
}

func (b *Bot) editApartmentDetail(ctx context.Context, chatID int64, msgID int, aptID int) {
	apt, err := b.aptRepo.GetByID(ctx, aptID)
	if err != nil {
		log.Printf("bot: GetByID(%d): %v", aptID, err)
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%s*\n", apt.Title))
	if apt.Address != "" {
		sb.WriteString(fmt.Sprintf("📍 %s\n", apt.Address))
	}
	sb.WriteString(fmt.Sprintf("\n🛏 Комнат: %d  |  👥 Макс. гостей: %d\n", apt.Rooms, apt.MaxGuests))
	sb.WriteString(fmt.Sprintf("💰 *%.0f грн / ночь*\n", apt.PricePerNight))
	if apt.Description != "" {
		sb.WriteString(fmt.Sprintf("\n%s\n", apt.Description))
	}
	if len(apt.Amenities) > 0 {
		sb.WriteString("\n✨ *Удобства:* " + strings.Join(apt.Amenities, ", ") + "\n")
	}

	btns := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📅 Забронировать на даты", fmt.Sprintf("%s%d", callbackBooking, aptID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💬 Задать вопрос", fmt.Sprintf("%s%d", callbackAskQuestion, aptID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к результатам", callbackSearch),
		),
	}

	// Загружаем фото квартиры
	photos, _ := b.photoRepo.GetByApartment(ctx, aptID)
	if len(photos) > 0 && b.uploadsRoot != "" {
		// Удаляем старое сообщение со списком
		_ = b.client.DeleteMessage(chatID, msgID)

		// Собираем абсолютные пути (до 10 фото — лимит Telegram)
		var filePaths []string
		for i, p := range photos {
			if i >= 10 {
				break
			}
			filePaths = append(filePaths, filepath.Join(b.uploadsRoot, filepath.FromSlash(p.FilePath)))
		}
		if err := b.client.SendMediaGroup(chatID, filePaths); err != nil {
			log.Printf("bot: SendMediaGroup apt=%d: %v", aptID, err)
		}
		// Текст с кнопкой отправляем новым сообщением
		_, _ = b.client.SendMessageWithKeyboardFull(chatID, sb.String(), btns)
	} else {
		// Фото нет — редактируем существующее сообщение
		_ = b.client.EditMessage(chatID, msgID, sb.String(), btns)
	}
}

// ─── Booking flow ─────────────────────────────────────────────────────────────

// startBookingCheckIn — показывает календарь выбора даты заезда.
func (b *Bot) startBookingCheckIn(chatID int64, msgID, aptID int) {
	b.mu.Lock()
	s := b.getOrCreateSession(chatID)
	s.bookingStep = 1
	s.bookingAptID = aptID
	s.bookingIn = time.Time{}
	s.bookingOut = time.Time{}
	b.mu.Unlock()

	ctx := context.Background()
	blockedDates, _ := b.bookingRepo.GetBlockedDates(ctx, aptID)

	now := time.Now()
	rows := buildCalendar(now.Year(), int(now.Month()), aptID, time.Time{}, now.Truncate(24*time.Hour), blockedDates)
	text := "📅 *Выберите дату заезда:*"
	_ = b.client.EditMessage(chatID, msgID, text, rows)
}

// handleCalNav — пользователь нажал < или > для смены месяца.
func (b *Bot) handleCalNav(chatID int64, msgID int, from *tgbotapi.User, payload string) {
	// payload: YYYY-MM:APTID
	parts := strings.SplitN(payload, ":", 2)
	if len(parts) != 2 {
		return
	}
	aptID, _ := strconv.Atoi(parts[1])
	var year, month int
	fmt.Sscanf(parts[0], "%d-%d", &year, &month)

	ctx := context.Background()
	blockedDates, _ := b.bookingRepo.GetBlockedDates(ctx, aptID)

	s := b.getOrCreateSession(chatID)
	var checkIn time.Time
	var minDay time.Time
	if s.bookingStep == 1 {
		minDay = time.Now().Truncate(24 * time.Hour)
		checkIn = time.Time{}
	} else {
		minDay = s.bookingIn.Add(24 * time.Hour)
		checkIn = s.bookingIn
	}

	rows := buildCalendar(year, month, aptID, checkIn, minDay, blockedDates)
	var text string
	if s.bookingStep == 2 {
		text = fmt.Sprintf("📅 *Выберите дату выезда:*\n\nЗаезд: *%s*", s.bookingIn.Format("02.01.2006"))
	} else {
		text = "📅 *Выберите дату заезда:*"
	}
	_ = b.client.EditMessage(chatID, msgID, text, rows)
}

// handleCalDay — пользователь выбрал конкретный день.
func (b *Bot) handleCalDay(ctx context.Context, chatID int64, msgID int, from *tgbotapi.User, payload string) {
	// payload: YYYY-MM-DD:APTID
	parts := strings.SplitN(payload, ":", 2)
	if len(parts) != 2 {
		return
	}
	aptID, _ := strconv.Atoi(parts[1])
	day, err := time.Parse("2006-01-02", parts[0])
	if err != nil {
		return
	}

	b.mu.Lock()
	s := b.getOrCreateSession(chatID)
	step := s.bookingStep
	b.mu.Unlock()

	if step == 1 {
		// Выбран заезд — показываем календарь выезда
		b.mu.Lock()
		s.bookingIn = day
		s.bookingStep = 2
		b.mu.Unlock()

		blockedDates, _ := b.bookingRepo.GetBlockedDates(ctx, aptID)
		minOut := day.Add(24 * time.Hour)
		rows := buildCalendar(day.Year(), int(day.Month()), aptID, day, minOut, blockedDates)
		text := fmt.Sprintf("📅 *Выберите дату выезда:*\n\nЗаезд: *%s*", day.Format("02.01.2006"))
		_ = b.client.EditMessage(chatID, msgID, text, rows)

	} else if step == 2 {
		// Выбран выезд — показываем подтверждение
		b.mu.Lock()
		s.bookingOut = day
		s.bookingStep = 3
		checkIn := s.bookingIn
		b.mu.Unlock()

		b.editBookingConfirm(ctx, chatID, msgID, aptID, checkIn, day)
	}
}

// editBookingConfirm — экран подтверждения бронирования.
func (b *Bot) editBookingConfirm(ctx context.Context, chatID int64, msgID, aptID int, checkIn, checkOut time.Time) {
	apt, err := b.aptRepo.GetByID(ctx, aptID)
	if err != nil || apt == nil {
		return
	}
	nights := int(checkOut.Sub(checkIn).Hours() / 24)
	total := float64(nights) * apt.PricePerNight

	text := fmt.Sprintf(
		"🏠 *%s*\n\n📅 Заезд: *%s*\n📅 Выезд: *%s*\n🌙 Ночей: *%d*\n💰 Итого: *%.0f грн*\n\nПодтвердить бронирование?",
		apt.Title,
		checkIn.Format("02.01.2006"),
		checkOut.Format("02.01.2006"),
		nights,
		total,
	)

	confirmPayload := fmt.Sprintf("%d:%s:%s", aptID,
		checkIn.Format("2006-01-02"), checkOut.Format("2006-01-02"))

	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить", callbackBookingConfirm+confirmPayload),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✖️ Отмена", fmt.Sprintf("%s%d", callbackBookingCancel, aptID)),
		),
	}
	_ = b.client.EditMessage(chatID, msgID, text, rows)
}

// confirmBooking — сохраняет бронирование в БД.
func (b *Bot) confirmBooking(ctx context.Context, chatID int64, msgID int, from *tgbotapi.User, payload string) {
	// payload: APTID:YYYY-MM-DD:YYYY-MM-DD
	parts := strings.SplitN(payload, ":", 3)
	if len(parts) != 3 {
		return
	}
	aptID, _ := strconv.Atoi(parts[0])
	checkIn, err1 := time.Parse("2006-01-02", parts[1])
	checkOut, err2 := time.Parse("2006-01-02", parts[2])
	if err1 != nil || err2 != nil {
		return
	}

	apt, err := b.aptRepo.GetByID(ctx, aptID)
	if err != nil || apt == nil {
		return
	}
	nights := int(checkOut.Sub(checkIn).Hours() / 24)
	total := float64(nights) * apt.PricePerNight

	username := ""
	firstName := ""
	if from != nil {
		username = from.UserName
		firstName = from.FirstName
	}
	clientID, err := b.bookingRepo.GetOrCreateUser(ctx, chatID, firstName, username)
	if err != nil {
		log.Printf("bot: confirmBooking GetOrCreateUser: %v", err)
		_ = b.client.EditMessage(chatID, msgID, "❌ Ошибка при создании бронирования. Попробуйте позже.", nil)
		return
	}

	booking, err := b.bookingRepo.Create(ctx, aptID, clientID, checkIn, checkOut, 1, total)
	if err != nil {
		log.Printf("bot: confirmBooking Create: %v", err)
		_ = b.client.EditMessage(chatID, msgID, "❌ Ошибка при создании бронирования. Попробуйте позже.", nil)
		return
	}

	b.cancelBookingFlow(chatID, aptID)

	text := fmt.Sprintf(
		"✅ *Заявка на бронирование принята!*\n\n🏠 %s\n📅 %s — %s\n🌙 Ночей: %d\n💰 Сумма: *%.0f грн*\n\n📋 Номер заявки: *#%d*\n\nС вами свяжутся для подтверждения.",
		apt.Title,
		checkIn.Format("02.01.2006"),
		checkOut.Format("02.01.2006"),
		nights,
		total,
		booking.ID,
	)

	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к квартире", fmt.Sprintf("%s%d", callbackApartment, aptID)),
		),
	}
	_ = b.client.EditMessage(chatID, msgID, text, rows)

	// Уведомляем владельца квартиры
	go b.notifyOwner(ctx, booking.ID, apt.Title, chatID, checkIn, checkOut, nights, total)
}

func (b *Bot) notifyOwner(ctx context.Context, bookingID int, aptTitle string, clientTgID int64, checkIn, checkOut time.Time, nights int, total float64) {
	ownerTgID, err := b.bookingRepo.GetOwnerTelegramIDByBooking(ctx, bookingID)
	if err != nil {
		log.Printf("notifyOwner: GetOwnerTelegramIDByBooking: %v", err)
		return
	}

	text := fmt.Sprintf(
		"🔔 *Новая заявка на бронирование #%d*\n\n🏠 %s\n📅 %s — %s\n🌙 Ночей: %d\n💰 Сумма: *%.0f грн*\n\n👤 Клиент: tg id %d",
		bookingID, aptTitle,
		checkIn.Format("02.01.2006"),
		checkOut.Format("02.01.2006"),
		nights, total, clientTgID,
	)

	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить", fmt.Sprintf("%s%d", callbackBookingApprove, bookingID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отклонить", fmt.Sprintf("%s%d", callbackBookingReject, bookingID)),
		),
	}
	_ = b.client.SendMessageWithKeyboard(ownerTgID, text, rows)
}

func (b *Bot) handleBookingApprove(ctx context.Context, chatID int64, msgID int, payload string) {
	bookingID, err := strconv.Atoi(payload)
	if err != nil {
		return
	}
	if err := b.bookingRepo.UpdateStatus(ctx, bookingID, "approved", ""); err != nil {
		log.Printf("handleBookingApprove: UpdateStatus: %v", err)
		return
	}
	_ = b.client.EditMessage(chatID, msgID, fmt.Sprintf("✅ Заявка *#%d* подтверждена.", bookingID), nil)

	// Уведомляем клиента
	booking, err := b.bookingRepo.GetByID(ctx, bookingID)
	if err != nil || booking == nil {
		return
	}
	clientTgID, err := b.bookingRepo.GetClientTelegramID(ctx, booking.ClientID)
	if err != nil {
		log.Printf("handleBookingApprove: GetClientTelegramID: %v", err)
		return
	}
	_ = b.client.SendMessage(clientTgID, fmt.Sprintf("🎉 Ваша заявка *#%d* подтверждена! Ждём вас!", bookingID))
}

func (b *Bot) handleBookingReject(ctx context.Context, chatID int64, msgID int, payload string) {
	bookingID, err := strconv.Atoi(payload)
	if err != nil {
		return
	}
	if err := b.bookingRepo.UpdateStatus(ctx, bookingID, "rejected", ""); err != nil {
		log.Printf("handleBookingReject: UpdateStatus: %v", err)
		return
	}
	_ = b.client.EditMessage(chatID, msgID, fmt.Sprintf("❌ Заявка *#%d* отклонена.", bookingID), nil)

	// Уведомляем клиента
	booking, err := b.bookingRepo.GetByID(ctx, bookingID)
	if err != nil || booking == nil {
		return
	}
	clientTgID, err := b.bookingRepo.GetClientTelegramID(ctx, booking.ClientID)
	if err != nil {
		log.Printf("handleBookingReject: GetClientTelegramID: %v", err)
		return
	}
	_ = b.client.SendMessage(clientTgID, fmt.Sprintf("😔 Ваша заявка *#%d* отклонена. Попробуйте выбрать другие даты.", bookingID))
}

func (b *Bot) cancelBookingFlow(chatID int64, aptID int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if s, ok := b.sessions[chatID]; ok {
		s.bookingStep = 0
		s.bookingAptID = 0
		s.bookingIn = time.Time{}
		s.bookingOut = time.Time{}
	}
}

func (b *Bot) getOrCreateSession(chatID int64) *session {
	if s, ok := b.sessions[chatID]; ok {
		return s
	}
	s := &session{filters: make(map[string]bool)}
	b.sessions[chatID] = s
	return s
}

// ─── Chat ─────────────────────────────────────────────────────────────────────

// handleAskQuestion — клиент нажал "Задать вопрос" на карточке квартиры.
func (b *Bot) handleAskQuestion(ctx context.Context, chatID int64, aptID int) {
	apt, err := b.aptRepo.GetByID(ctx, aptID)
	if err != nil || apt == nil {
		_ = b.client.SendMessage(chatID, "❌ Не удалось найти квартиру.")
		return
	}
	owner, err := b.userRepo.GetByID(ctx, apt.OwnerID)
	if err != nil || owner == nil {
		_ = b.client.SendMessage(chatID, "❌ Арендодатель не найден.")
		return
	}

	b.mu.Lock()
	s := b.getOrCreateSession(chatID)
	s.chatAptID = aptID
	s.chatOwnerTgID = owner.TelegramID
	s.chatAptTitle = apt.Title
	b.mu.Unlock()

	_ = b.client.SendMessage(chatID,
		"💬 Вы в режиме чата с арендодателем квартиры «"+apt.Title+"».\n\nНапишите ваш вопрос — я перешлю его.\nДля выхода из чата — /stopchat")
}

// handleChatMessage — обрабатывает текстовые сообщения (не команды).
func (b *Bot) handleChatMessage(ctx context.Context, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	text := msg.Text
	if text == "" {
		return
	}

	s := b.getSession(chatID)

	// Клиент пишет арендодателю
	if s.chatAptID > 0 {
		senderName := msg.From.FirstName
		if msg.From.UserName != "" {
			senderName += " (@" + msg.From.UserName + ")"
		}
		fwdText := fmt.Sprintf("📨 *Вопрос от клиента* %s\n🏠 Квартира: %s\n\n%s",
			escapeMarkdownV2(senderName),
			escapeMarkdownV2(s.chatAptTitle),
			escapeMarkdownV2(text),
		)
		rows := [][]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("💬 Ответить клиенту", fmt.Sprintf("%s%d", callbackChatReply, chatID)),
			),
		}
		_ = b.client.SendMessageWithKeyboard(s.chatOwnerTgID, fwdText, rows)
		_ = b.client.SendMessage(chatID, "✅ Ваш вопрос отправлен арендодателю.")
		return
	}

	// Арендодатель отвечает клиенту
	if s.replyToClientID != 0 {
		fwdText := fmt.Sprintf("📨 *Ответ арендодателя:*\n\n%s",
			escapeMarkdownV2(text),
		)
		_ = b.client.sendMarkdownV2(s.replyToClientID, fwdText)

		b.mu.Lock()
		s.replyToClientID = 0
		b.mu.Unlock()

		_ = b.client.SendMessage(chatID, "✅ Ответ отправлен клиенту.")
		return
	}
}

// handleChatReplyMode — арендодатель нажал "Ответить клиенту".
func (b *Bot) handleChatReplyMode(ownerChatID int64, clientChatID int64) {
	b.mu.Lock()
	s := b.getOrCreateSession(ownerChatID)
	s.replyToClientID = clientChatID
	b.mu.Unlock()

	_ = b.client.SendMessage(ownerChatID, "✏️ Режим ответа активирован. Напишите следующее сообщение — оно будет передано клиенту.")
}

// cmdStopChat — выход из режима чата.
func (b *Bot) cmdStopChat(msg *tgbotapi.Message) {
	b.mu.Lock()
	s := b.getOrCreateSession(msg.Chat.ID)
	s.chatAptID = 0
	s.chatOwnerTgID = 0
	s.chatAptTitle = ""
	s.replyToClientID = 0
	b.mu.Unlock()

	_ = b.client.SendMessage(msg.Chat.ID, "🔕 Режим чата завершён.")
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func escapeMarkdownV2(s string) string {
	replacer := strings.NewReplacer(
		"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]",
		"(", "\\(", ")", "\\)", "~", "\\~", "`", "\\`",
		">", "\\>", "#", "\\#", "+", "\\+", "-", "\\-",
		"=", "\\=", "|", "\\|", "{", "\\{", "}", "\\}",
		".", "\\.", "!", "\\!",
	)
	return replacer.Replace(s)
}
