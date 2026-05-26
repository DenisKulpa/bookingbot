package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/DenisKulpa/bookingbot/internal/repository"
)

const (
	callbackDistrict   = "district:"
	callbackApartment  = "apartment:"
	callbackBack       = "back:districts"
	callbackFilterCat  = "filter_cat:"
	callbackFilter     = "filter:"
	callbackBackFilters = "back:filters"
)

// ─── Filter data ─────────────────────────────────────────────────────────────

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

// ─── Bot ──────────────────────────────────────────────────────────────────────

type Bot struct {
	client   *Client
	zoneRepo *repository.ZoneRepository
	aptRepo  *repository.ApartmentRepository
}

func NewBot(client *Client, zoneRepo *repository.ZoneRepository, aptRepo *repository.ApartmentRepository) *Bot {
	return &Bot{
		client:   client,
		zoneRepo: zoneRepo,
		aptRepo:  aptRepo,
	}
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

// ─── Message commands ─────────────────────────────────────────────────────────

func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	if !msg.IsCommand() {
		return
	}

	switch msg.Command() {
	case "start":
		b.cmdStart(ctx, msg)
	case "search":
		b.cmdSearch(ctx, msg)
	case "help":
		b.cmdHelp(ctx, msg)
	default:
		_ = b.client.SendMessage(msg.Chat.ID, "Неизвестная команда. Используйте /help для списка команд.")
	}
}

func (b *Bot) cmdStart(ctx context.Context, msg *tgbotapi.Message) {
	name := msg.From.FirstName
	if name == "" {
		name = msg.From.UserName
	}
	text := fmt.Sprintf(
		"👋 Привет, *%s*\\!\n\n"+
			"Я помогу найти и забронировать квартиру в Одессе\\.\n\n"+
			"• /search — поиск жилья по фильтрам\n"+
			"• /help — список команд",
		escapeMarkdownV2(name),
	)
	_ = b.client.sendMarkdownV2(msg.Chat.ID, text)
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

// ─── Callbacks ────────────────────────────────────────────────────────────────

func (b *Bot) handleCallback(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	_ = b.client.AnswerCallbackQuery(cq.ID, "", false)

	data := cq.Data
	chatID := cq.Message.Chat.ID
	msgID := cq.Message.MessageID

	switch {
	// ── New filter flow ──
	case data == callbackBackFilters:
		b.editFilterList(chatID, msgID)

	case strings.HasPrefix(data, callbackFilterCat):
		catCode := strings.TrimPrefix(data, callbackFilterCat)
		b.editFilterOptions(chatID, msgID, catCode)

	case strings.HasPrefix(data, callbackFilter):
		_ = b.client.AnswerCallbackQuery(cq.ID, "✅ Выбрано", false)

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
	}
}

// ─── Filter screens ───────────────────────────────────────────────────────────

func (b *Bot) sendFilterList(chatID int64) {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, cat := range filterCategories {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(cat.Label, callbackFilterCat+cat.Code),
		))
	}
	_ = b.client.SendMessageWithKeyboard(chatID, "🏖 *Поиск жилья в Аркадии по фильтрам*\n\nВыберите категорию:", rows)
}

func (b *Bot) editFilterList(chatID int64, msgID int) {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, cat := range filterCategories {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(cat.Label, callbackFilterCat+cat.Code),
		))
	}
	_ = b.client.EditMessage(chatID, msgID, "🏖 *Поиск жилья в Аркадии по фильтрам*\n\nВыберите категорию:", rows)
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
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, opt := range cat.Options {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(opt.Label, callbackFilter+opt.Code),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к фильтрам", callbackBackFilters),
	))
	_ = b.client.EditMessage(chatID, msgID, cat.Label, rows)
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

	var districtID int
	if apt.ZoneID != nil {
		districtID = *apt.ZoneID
	}

	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к квартирам", callbackDistrict+fmt.Sprint(districtID)),
		),
	}

	_ = b.client.EditMessage(chatID, msgID, sb.String(), rows)
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
