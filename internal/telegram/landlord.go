package telegram

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/DenisKulpa/bookingbot/internal/model"
)

// ─── Меню арендодателя ────────────────────────────────────────────────────────

func (b *Bot) sendLandlordMenu(chatID int64, name string) {
	text := fmt.Sprintf("👋 Добро пожаловать, *%s*!\n\n🏠 Панель арендодателя", name)
	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Мои квартиры", callbackMyApartments),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить квартиру", callbackAddApartment),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Входящие брони", callbackIncomingBookings),
		),
	}
	_ = b.client.SendMessageWithKeyboard(chatID, text, rows)
}

// ─── Мои квартиры ─────────────────────────────────────────────────────────────

func (b *Bot) sendMyApartments(ctx context.Context, chatID int64, msgID int) {
	user, err := b.userRepo.GetByTelegramID(ctx, chatID)
	if err != nil || user == nil {
		_ = b.client.SendMessage(chatID, "❌ Пользователь не найден.")
		return
	}
	apts, err := b.aptRepo.GetByOwner(ctx, user.ID)
	if err != nil {
		log.Printf("sendMyApartments: %v", err)
		_ = b.client.SendMessage(chatID, "❌ Ошибка загрузки квартир.")
		return
	}

	if len(apts) == 0 {
		rows := [][]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➕ Добавить квартиру", callbackAddApartment),
			),
		}
		_ = b.client.EditMessage(chatID, msgID, "🏠 *Мои квартиры*\n\nУ вас пока нет квартир.", rows)
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, apt := range apts {
		status := "✅"
		if !apt.IsAvailable {
			status = "❌"
		}
		label := fmt.Sprintf("%s %s — %.0f грн/ночь", status, apt.Title, apt.PricePerNight)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, callbackApartment+strconv.Itoa(apt.ID)),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("➕ Добавить квартиру", callbackAddApartment),
	))
	_ = b.client.EditMessage(chatID, msgID, fmt.Sprintf("🏠 *Мои квартиры* (%d)", len(apts)), rows)
}

// ─── Входящие брони ───────────────────────────────────────────────────────────

func (b *Bot) sendIncomingBookings(ctx context.Context, chatID int64, msgID int) {
	user, err := b.userRepo.GetByTelegramID(ctx, chatID)
	if err != nil || user == nil {
		_ = b.client.SendMessage(chatID, "❌ Пользователь не найден.")
		return
	}
	bookings, err := b.bookingRepo.GetPendingByOwner(ctx, user.ID)
	if err != nil {
		log.Printf("sendIncomingBookings: %v", err)
		_ = b.client.SendMessage(chatID, "❌ Ошибка загрузки броней.")
		return
	}

	if len(bookings) == 0 {
		_ = b.client.EditMessage(chatID, msgID, "📋 *Входящие брони*\n\nНет новых заявок.", nil)
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 *Входящие брони* (%d)\n\n", len(bookings)))
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, pb := range bookings {
		nights := int(pb.CheckOut.Sub(pb.CheckIn).Hours() / 24)
		clientLabel := pb.ClientName
		if pb.ClientUsername != "" {
			clientLabel += " (@" + pb.ClientUsername + ")"
		}
		sb.WriteString(fmt.Sprintf("📌 *#%d* — %s\n📅 %s – %s (%d ноч.)\n💰 %.0f грн\n👤 %s\n\n",
			pb.ID, pb.ApartmentTitle,
			pb.CheckIn.Format("02.01"), pb.CheckOut.Format("02.01"), nights,
			pb.TotalPrice, clientLabel,
		))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("✅ #%d", pb.ID), callbackBookingApprove+strconv.Itoa(pb.ID)),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("❌ #%d", pb.ID), callbackBookingReject+strconv.Itoa(pb.ID)),
		))
	}
	_ = b.client.EditMessage(chatID, msgID, sb.String(), rows)
}

// ─── Визард добавления квартиры ───────────────────────────────────────────────

const (
	wStepTitle   = 1
	wStepDesc    = 2
	wStepAddr    = 3
	wStepZone    = 4
	wStepRooms   = 5
	wStepPrice   = 6
	wStepPhotos  = 7
	wStepFilters = 8
	wStepConfirm = 9
)

func (b *Bot) startAddApartment(ctx context.Context, chatID int64) {
	b.mu.Lock()
	s := b.getOrCreateSession(chatID)
	s.addAptStep = wStepTitle
	s.addAptTitle = ""
	s.addAptDesc = ""
	s.addAptAddr = ""
	s.addAptZoneID = 0
	s.addAptRooms = 1
	s.addAptPrice = 0
	s.addAptPhotoIDs = nil
	s.addAptFilters = make(map[string]bool)
	b.mu.Unlock()

	_ = b.client.SendMessage(chatID, "🏠 *Добавление квартиры*\n\nШаг 1/8 — Введите *название* квартиры:")
}

// wizardHandleText — обработчик текстового ввода на каждом шаге визарда.
func (b *Bot) wizardHandleText(ctx context.Context, chatID int64, text string) {
	s := b.getSession(chatID)

	switch s.addAptStep {
	case wStepTitle:
		b.mu.Lock()
		if text != "-" {
			s.addAptTitle = text
		}
		s.addAptStep = wStepDesc
		b.mu.Unlock()
		_ = b.client.SendMessage(chatID, "Шаг 2/8 — Введите *описание* квартиры (или отправьте \"-\" чтобы пропустить):")

	case wStepDesc:
		b.mu.Lock()
		if text != "-" {
			s.addAptDesc = text
		}
		s.addAptStep = wStepAddr
		b.mu.Unlock()
		_ = b.client.SendMessage(chatID, "Шаг 3/8 — Введите *адрес* квартиры (или \"-\" чтобы пропустить):")

	case wStepAddr:
		b.mu.Lock()
		if text != "-" {
			s.addAptAddr = text
		}
		s.addAptStep = wStepZone
		b.mu.Unlock()
		b.wizardShowZones(ctx, chatID)

	case wStepRooms:
		if text == "-" {
			b.mu.Lock()
			s.addAptStep = wStepPrice
			b.mu.Unlock()
			_ = b.client.SendMessage(chatID, "Шаг 6/8 — Введите *цену за ночь* в гривнах (только число, или \"-\" чтобы оставить):")
			return
		}
		rooms, err := strconv.Atoi(strings.TrimSpace(text))
		if err != nil || rooms < 1 || rooms > 20 {
			_ = b.client.SendMessage(chatID, "⚠️ Введите число комнат (от 1 до 20) или \"-\" чтобы пропустить:")
			return
		}
		b.mu.Lock()
		s.addAptRooms = rooms
		s.addAptStep = wStepPrice
		b.mu.Unlock()
		_ = b.client.SendMessage(chatID, "Шаг 6/8 — Введите *цену за ночь* в гривнах (только число):")

	case wStepPrice:
		if text == "-" {
			b.mu.Lock()
			s.addAptStep = wStepPhotos
			b.mu.Unlock()
			b.wizardAskPhotos(chatID)
			return
		}
		priceStr := strings.ReplaceAll(strings.TrimSpace(text), ",", ".")
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil || price <= 0 {
			_ = b.client.SendMessage(chatID, "⚠️ Введите корректную цену (например: 1500) или \"-\" чтобы оставить:")
			return
		}
		b.mu.Lock()
		s.addAptPrice = price
		s.addAptStep = wStepPhotos
		b.mu.Unlock()
		b.wizardAskPhotos(chatID)
	}
}

// wizardShowZones — показывает список зон для выбора.
func (b *Bot) wizardShowZones(ctx context.Context, chatID int64) {
	zones, err := b.zoneRepo.GetSubzonesFlat(ctx, 3) // parent_id=3 Аркадия
	if err != nil {
		log.Printf("wizardShowZones: %v", err)
		_ = b.client.SendMessage(chatID, "❌ Ошибка загрузки зон.")
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, z := range zones {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(z.Name, callbackWizardZone+strconv.Itoa(z.ID)),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⏭ Пропустить", callbackWizardZone+"0"),
	))
	_ = b.client.SendMessageWithKeyboard(chatID, "Шаг 4/8 — Выберите *зону* в Аркадии:", rows)
}

// wizardSelectZone — пользователь выбрал зону.
func (b *Bot) wizardSelectZone(ctx context.Context, chatID int64, msgID int, zoneID int) {
	b.mu.Lock()
	s := b.getOrCreateSession(chatID)
	s.addAptZoneID = zoneID
	s.addAptStep = wStepRooms
	b.mu.Unlock()

	_ = b.client.EditMessage(chatID, msgID, "✅ Зона выбрана!", nil)
	_ = b.client.SendMessage(chatID, "Шаг 5/8 — Введите *количество комнат* (1, 2, 3 и т.д.):")
}

// wizardAskPhotos — предлагает загрузить фото.
func (b *Bot) wizardAskPhotos(chatID int64) {
	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Готово (фото добавлены)", callbackWizardDonePhotos),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏭ Пропустить фото", callbackWizardDonePhotos),
		),
	}
	_ = b.client.SendMessageWithKeyboard(chatID,
		"Шаг 7/8 — Отправьте *фотографии* квартиры (по одной или несколько).\nКогда закончите — нажмите «Готово».", rows)
}

// handleWizardPhoto — получено фото в режиме визарда.
func (b *Bot) handleWizardPhoto(ctx context.Context, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	s := b.getSession(chatID)
	if s.addAptStep != wStepPhotos {
		return
	}
	photos := msg.Photo
	if len(photos) == 0 {
		return
	}
	// Берём самое большое фото
	best := photos[len(photos)-1]
	b.mu.Lock()
	s.addAptPhotoIDs = append(s.addAptPhotoIDs, best.FileID)
	count := len(s.addAptPhotoIDs)
	b.mu.Unlock()

	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Готово (фото добавлены)", callbackWizardDonePhotos),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏭ Пропустить фото", callbackWizardDonePhotos),
		),
	}
	_ = b.client.SendMessageWithKeyboard(chatID,
		fmt.Sprintf("📷 Фото добавлено (%d шт.). Продолжайте отправлять или нажмите «Готово».", count), rows)
}

// wizardDonePhotos — фото загружены, переходим к фильтрам.
func (b *Bot) wizardDonePhotos(chatID int64, msgID int) {
	b.mu.Lock()
	s := b.getOrCreateSession(chatID)
	s.addAptStep = wStepFilters
	b.mu.Unlock()

	_ = b.client.EditMessage(chatID, msgID, "✅ Фото приняты!", nil)
	b.wizardShowFilterList(chatID)
}

// wizardShowFilterList — главный экран выбора фильтров для визарда.
func (b *Bot) wizardShowFilterList(chatID int64) {
	s := b.getSession(chatID)
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, cat := range b.filterCats {
		hasActive := false
		for _, opt := range cat.Options {
			if s.addAptFilters[opt.Code] {
				hasActive = true
				break
			}
		}
		label := cat.Label
		if hasActive {
			label += " ✅"
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, callbackWizardFilterCat+cat.Code),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅ Готово (перейти к подтверждению)", callbackWizardDoneFilters),
	))
	_ = b.client.SendMessageWithKeyboard(chatID, "Шаг 8/8 — Выберите *удобства и фильтры* квартиры:", rows)
}

// wizardShowFilterCat — показывает опции одной категории фильтров.
func (b *Bot) wizardShowFilterCat(chatID int64, msgID int, catCode string) {
	var cat *model.FilterCategory
	for i := range b.filterCats {
		if b.filterCats[i].Code == catCode {
			cat = &b.filterCats[i]
			break
		}
	}
	if cat == nil {
		return
	}
	s := b.getSession(chatID)
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, opt := range cat.Options {
		label := "➖  " + opt.Label
		if s.addAptFilters[opt.Code] {
			label = "✔️  " + opt.Label
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, callbackWizardToggle+opt.Code),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к категориям", callbackWizardBackFilters),
	))
	_ = b.client.EditMessage(chatID, msgID, cat.Label+":", rows)
}

// wizardToggleFilter — переключает фильтр в визарде.
func (b *Bot) wizardToggleFilter(chatID int64, code string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := b.getOrCreateSession(chatID)
	if s.addAptFilters == nil {
		s.addAptFilters = make(map[string]bool)
	}
	s.addAptFilters[code] = !s.addAptFilters[code]
}

// wizardEditFilterList — возврат к списку фильтров (редактирует сообщение).
func (b *Bot) wizardEditFilterList(chatID int64, msgID int) {
	s := b.getSession(chatID)
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, cat := range b.filterCats {
		hasActive := false
		for _, opt := range cat.Options {
			if s.addAptFilters[opt.Code] {
				hasActive = true
				break
			}
		}
		label := cat.Label
		if hasActive {
			label += " ✅"
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, callbackWizardFilterCat+cat.Code),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅ Готово (перейти к подтверждению)", callbackWizardDoneFilters),
	))
	_ = b.client.EditMessage(chatID, msgID, "Шаг 8/8 — Выберите *удобства и фильтры* квартиры:", rows)
}

// wizardShowConfirm — экран подтверждения перед сохранением.
func (b *Bot) wizardShowConfirm(ctx context.Context, chatID int64, msgID int) {
	s := b.getSession(chatID)

	_ = b.client.EditMessage(chatID, msgID, "✅ Фильтры выбраны!", nil)

	zoneLabel := "не указана"
	if s.addAptZoneID > 0 {
		zones, err := b.zoneRepo.GetSubzonesFlat(ctx, 3)
		if err == nil {
			for _, z := range zones {
				if z.ID == s.addAptZoneID {
					zoneLabel = z.Name
					break
				}
			}
		}
	}

	var filterLabels []string
	for _, cat := range b.filterCats {
		for _, opt := range cat.Options {
			if s.addAptFilters[opt.Code] {
				filterLabels = append(filterLabels, opt.Label)
			}
		}
	}
	filtersText := "нет"
	if len(filterLabels) > 0 {
		filtersText = strings.Join(filterLabels, ", ")
	}

	text := fmt.Sprintf(
		"📋 *Проверьте данные квартиры:*\n\n"+
			"🏷 Название: *%s*\n"+
			"📝 Описание: %s\n"+
			"📍 Адрес: %s\n"+
			"🗺 Зона: %s\n"+
			"🛏 Комнат: *%d*\n"+
			"💰 Цена: *%.0f грн/ночь*\n"+
			"📷 Фото: %d шт.\n"+
			"✨ Фильтры: %s",
		s.addAptTitle,
		strDefault(s.addAptDesc, "не указано"),
		strDefault(s.addAptAddr, "не указан"),
		zoneLabel,
		s.addAptRooms,
		s.addAptPrice,
		len(s.addAptPhotoIDs),
		filtersText,
	)

	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Сохранить квартиру", callbackWizardConfirm),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", callbackWizardCancel),
		),
	}
	_ = b.client.SendMessageWithKeyboard(chatID, text, rows)
}

// wizardSaveApartment — финальное сохранение в БД (создание или обновление).
func (b *Bot) wizardSaveApartment(ctx context.Context, chatID int64, msgID int, from *tgbotapi.User) {
	s := b.getSession(chatID)

	user, err := b.userRepo.GetByTelegramID(ctx, chatID)
	if err != nil || user == nil {
		_ = b.client.EditMessage(chatID, msgID, "❌ Ошибка: пользователь не найден.", nil)
		return
	}

	var zoneIDPtr *int
	if s.addAptZoneID > 0 {
		zid := s.addAptZoneID
		zoneIDPtr = &zid
	}

	// Извлекаем тип квартиры из фильтров (type_studio → studio, type_1room → 1room, ...)
	aptType := ""
	for code, active := range s.addAptFilters {
		if active && strings.HasPrefix(code, "type_") {
			aptType = strings.TrimPrefix(code, "type_")
			break
		}
	}

	var aptID int
	isEdit := s.editAptID > 0

	if isEdit {
		// Обновление существующей квартиры
		aptID = s.editAptID
		_, err = b.aptRepo.UpdateFull(ctx, aptID, zoneIDPtr, s.addAptTitle, s.addAptDesc, s.addAptAddr,
			aptType, s.addAptRooms, s.addAptPrice)
		if err != nil {
			log.Printf("wizardSaveApartment UpdateFull: %v", err)
			_ = b.client.EditMessage(chatID, msgID, "❌ Ошибка при обновлении квартиры.", nil)
			return
		}
		// Удаляем старые фильтры и фото — будут пересозданы ниже
		_ = b.aptRepo.ClearFilters(ctx, aptID)
	} else {
		// Создание новой
		apt, err := b.aptRepo.Create(ctx,
			user.ID, zoneIDPtr,
			s.addAptTitle, s.addAptDesc, s.addAptAddr, aptType,
			s.addAptRooms, s.addAptRooms*2,
			s.addAptPrice,
		)
		if err != nil {
			log.Printf("wizardSaveApartment Create: %v", err)
			_ = b.client.EditMessage(chatID, msgID, "❌ Ошибка при создании квартиры.", nil)
			return
		}
		aptID = apt.ID
	}

	// Сохраняем фото (только новые, загруженные в этом визарде)
	aptDir := filepath.Join(b.uploadsRoot, "uploads", "apartments", strconv.Itoa(aptID))
	_ = os.MkdirAll(aptDir, 0755)
	for i, fileID := range s.addAptPhotoIDs {
		ext, err := b.client.GetFileExt(fileID)
		if err != nil {
			log.Printf("wizardSaveApartment GetFileExt[%d]: %v", i, err)
			ext = ".jpg"
		}
		if ext == "" {
			ext = ".jpg"
		}
		filename := fmt.Sprintf("%d%s", i+1, ext)
		destPath := filepath.Join(aptDir, filename)

		if err := b.client.DownloadPhoto(fileID, destPath); err != nil {
			log.Printf("wizardSaveApartment DownloadPhoto[%d]: %v", i, err)
			continue
		}

		relPath := filepath.ToSlash(filepath.Join("uploads", "apartments", strconv.Itoa(aptID), filename))
		if _, err := b.photoRepo.Add(ctx, aptID, relPath, "", i+1); err != nil {
			log.Printf("wizardSaveApartment photoRepo.Add[%d]: %v", i, err)
		}
	}

	// Сохраняем фильтры (удобства + зона + тип — всё в apartment_filters)
	var filterCodes []string
	for code, active := range s.addAptFilters {
		if active {
			filterCodes = append(filterCodes, code)
		}
	}
	if s.addAptZoneID > 0 {
		zoneCodes, err := b.zoneRepo.GetFilterCodes(ctx, s.addAptZoneID)
		if err != nil {
			log.Printf("wizardSaveApartment GetFilterCodes: %v", err)
		} else {
			filterCodes = append(filterCodes, zoneCodes...)
		}
	}
	if aptType != "" {
		typeCode := "type_" + aptType
		found := false
		for _, c := range filterCodes {
			if c == typeCode {
				found = true
				break
			}
		}
		if !found {
			filterCodes = append(filterCodes, typeCode)
		}
	}
	if len(filterCodes) > 0 {
		if err := b.aptRepo.AddFilters(ctx, aptID, filterCodes); err != nil {
			log.Printf("wizardSaveApartment AddFilters: %v", err)
		}
	}

	// Сбрасываем визард
	b.mu.Lock()
	s2 := b.getOrCreateSession(chatID)
	s2.addAptStep = 0
	s2.editAptID = 0
	s2.addAptPhotoIDs = nil
	s2.addAptFilters = nil
	b.mu.Unlock()

	if isEdit {
		_ = b.client.EditMessage(chatID, msgID,
			fmt.Sprintf("✅ *Квартира #%d обновлена!*\n\n🏠 %s\n💰 %.0f грн/ночь",
				aptID, s.addAptTitle, s.addAptPrice), nil)
	} else {
		_ = b.client.EditMessage(chatID, msgID,
			fmt.Sprintf("🎉 *Квартира #%d успешно добавлена!*\n\n🏠 %s\n💰 %.0f грн/ночь",
				aptID, s.addAptTitle, s.addAptPrice), nil)
	}
	b.sendLandlordMenu(chatID, from.FirstName)
}

// wizardCancel — отмена визарда.
func (b *Bot) wizardCancel(chatID int64, msgID int) {
	b.mu.Lock()
	s := b.getOrCreateSession(chatID)
	s.addAptStep = 0
	s.editAptID = 0
	b.mu.Unlock()

	_ = b.client.EditMessage(chatID, msgID, "❌ Добавление квартиры отменено.", nil)
}

// ─── Редактирование квартиры (тот же визард, что и создание) ─────────────────

func (b *Bot) startEditApartment(ctx context.Context, chatID int64, msgID int, aptID int) {
	apt, err := b.aptRepo.GetByID(ctx, aptID)
	if err != nil || apt == nil {
		_ = b.client.EditMessage(chatID, msgID, "❌ Квартира не найдена.", nil)
		return
	}

	// Загружаем текущие фильтры
	existingFilters := make(map[string]bool)
	filterCodes, _ := b.aptRepo.GetFilterCodes(ctx, aptID)
	for _, code := range filterCodes {
		existingFilters[code] = true
	}

	zoneID := 0
	if apt.ZoneID != nil {
		zoneID = *apt.ZoneID
	}

	b.mu.Lock()
	s := b.getOrCreateSession(chatID)
	s.editAptID = aptID
	s.addAptStep = wStepTitle
	s.addAptTitle = apt.Title
	s.addAptDesc = apt.Description
	s.addAptAddr = apt.Address
	s.addAptZoneID = zoneID
	s.addAptRooms = apt.Rooms
	s.addAptPrice = apt.PricePerNight
	s.addAptPhotoIDs = nil
	s.addAptFilters = existingFilters
	b.mu.Unlock()

	_ = b.client.EditMessage(chatID, msgID,
		fmt.Sprintf("✏️ *Редактирование квартиры #%d*\n\n"+
			"Шаг 1/8 — Введите *название*\nТекущее: *%s*\n(или \"-\" чтобы оставить)",
			aptID, apt.Title), nil)
}

// handleDeletePhoto — удаление фото при редактировании (из списка существующих фото).
func (b *Bot) handleDeletePhoto(ctx context.Context, chatID int64, msgID int, photoID int) {
	if err := b.photoRepo.Delete(ctx, photoID); err != nil {
		log.Printf("handleDeletePhoto: %v", err)
	}
	_ = b.client.AnswerCallbackQuery("", "", false)
	_ = b.client.DeleteMessage(chatID, msgID)
}

func (b *Bot) cancelEdit(chatID int64) {
	b.mu.Lock()
	s := b.getOrCreateSession(chatID)
	s.editAptID = 0
	s.editAptStep = 0
	b.mu.Unlock()
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func strDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
