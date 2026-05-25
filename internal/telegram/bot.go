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
	callbackDistrict  = "district:"
	callbackApartment = "apartment:"
	callbackBack      = "back:districts"
)

type Bot struct {
	client    *Client
	zoneRepo  *repository.ZoneRepository
	aptRepo   *repository.ApartmentRepository
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

// ─── Message commands ────────────────────────────────────────────────────────

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
			"• /search — просмотр районов и квартир\n"+
			"• /help — список команд",
		escapeMarkdownV2(name),
	)
	_ = b.client.sendMarkdownV2(msg.Chat.ID, text)
}

func (b *Bot) cmdHelp(_ context.Context, msg *tgbotapi.Message) {
	text := "*Список команд:*\n\n" +
		"/start — приветствие\n" +
		"/search — выбрать район и посмотреть квартиры\n" +
		"/help — это сообщение"
	_ = b.client.sendMarkdownV2(msg.Chat.ID, text)
}

func (b *Bot) cmdSearch(ctx context.Context, msg *tgbotapi.Message) {
	b.sendDistrictList(ctx, msg.Chat.ID, 0)
}

// ─── Callbacks ───────────────────────────────────────────────────────────────

func (b *Bot) handleCallback(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	_ = b.client.AnswerCallbackQuery(cq.ID, "", false)

	data := cq.Data
	chatID := cq.Message.Chat.ID
	msgID := cq.Message.MessageID

	switch {
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

// ─── Screens ─────────────────────────────────────────────────────────────────

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

	// Apartments button
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
