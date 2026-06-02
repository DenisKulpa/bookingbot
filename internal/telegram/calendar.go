package telegram

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Callback-форматы:
//   cal_nav:YYYY-MM:APTID        — переключить месяц
//   cal_day:YYYY-MM-DD:APTID     — выбрать день
//   cal_ignore                   — кнопка-заглушка (дни заголовка, пустые ячейки)

const (
	callbackCalNav    = "cal_nav:"
	callbackCalDay    = "cal_day:"
	callbackCalIgnore = "cal_ignore"
)

var weekDays = []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}

var monthNames = []string{
	"", "Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
	"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
}

// buildCalendar строит inline-клавиатуру календаря.
// aptID — id квартиры, нужен в callback для сохранения контекста.
// checkIn — если уже выбран заезд (ненулевой), подсвечиваем его.
// minDay — минимальная доступная дата (сегодня или день после checkIn).
// blockedDates — map[YYYY-MM-DD]true из БД (bookings + availability).
func buildCalendar(year, month, aptID int, checkIn time.Time, minDay time.Time, blockedDates map[string]bool) [][]tgbotapi.InlineKeyboardButton {
	var rows [][]tgbotapi.InlineKeyboardButton

	// ── Заголовок: < Месяц Год > ──────────────────────────────────────────
	monthKey := fmt.Sprintf("%04d-%02d", year, month)
	prevYear, prevMonth := year, month-1
	if prevMonth == 0 {
		prevMonth = 12
		prevYear--
	}
	nextYear, nextMonth := year, month+1
	if nextMonth == 13 {
		nextMonth = 1
		nextYear++
	}
	prevKey := fmt.Sprintf("%04d-%02d:%d", prevYear, prevMonth, aptID)
	nextKey := fmt.Sprintf("%04d-%02d:%d", nextYear, nextMonth, aptID)

	// Не пускать назад раньше текущего месяца
	today := time.Now()
	prevDisabled := (prevYear < today.Year()) || (prevYear == today.Year() && prevMonth < int(today.Month()))

	headerRow := tgbotapi.NewInlineKeyboardRow(
		func() tgbotapi.InlineKeyboardButton {
			if prevDisabled {
				return tgbotapi.NewInlineKeyboardButtonData(" ", callbackCalIgnore)
			}
			return tgbotapi.NewInlineKeyboardButtonData("◀", callbackCalNav+prevKey)
		}(),
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %d", monthNames[month], year),
			callbackCalIgnore,
		),
		tgbotapi.NewInlineKeyboardButtonData("▶", callbackCalNav+nextKey),
	)
	rows = append(rows, headerRow)

	// ── Дни недели ────────────────────────────────────────────────────────
	var wdRow []tgbotapi.InlineKeyboardButton
	for _, d := range weekDays {
		wdRow = append(wdRow, tgbotapi.NewInlineKeyboardButtonData(d, callbackCalIgnore))
	}
	rows = append(rows, wdRow)

	// ── Дни месяца ────────────────────────────────────────────────────────
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	// Понедельник = 1, воскресенье = 7
	startWeekday := int(firstDay.Weekday())
	if startWeekday == 0 {
		startWeekday = 7
	}
	daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()

	cell := 1
	var row []tgbotapi.InlineKeyboardButton

	// Пустые ячейки до первого числа
	for i := 1; i < startWeekday; i++ {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", callbackCalIgnore))
		cell++
	}

	for day := 1; day <= daysInMonth; day++ {
		d := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		label := fmt.Sprintf("%d", day)

		// Отметить выбранный заезд
		if !checkIn.IsZero() && d.Equal(checkIn) {
			label = "✅" + label
		}

		dayKey := fmt.Sprintf("%s-%02d:%d", monthKey, day, aptID)
		isBlocked := blockedDates[d.Format("2006-01-02")]

		var btn tgbotapi.InlineKeyboardButton
		if d.Before(minDay) || isBlocked {
			// Недоступный или занятый день
			if isBlocked && !d.Before(minDay) {
				btn = tgbotapi.NewInlineKeyboardButtonData("✖", callbackCalIgnore)
			} else {
				btn = tgbotapi.NewInlineKeyboardButtonData("·", callbackCalIgnore)
			}
		} else {
			btn = tgbotapi.NewInlineKeyboardButtonData(label, callbackCalDay+dayKey)
		}

		row = append(row, btn)
		cell++
		if cell > 7 {
			rows = append(rows, row)
			row = nil
			cell = 1
		}
	}
	// Последняя строка
	if len(row) > 0 {
		for len(row) < 7 {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", callbackCalIgnore))
		}
		rows = append(rows, row)
	}

	// ── Легенда ───────────────────────────────────────────────────────────
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("· прошедшие", callbackCalIgnore),
		tgbotapi.NewInlineKeyboardButtonData("✖ занято", callbackCalIgnore),
		tgbotapi.NewInlineKeyboardButtonData("✅ заезд", callbackCalIgnore),
	))

	// ── Кнопка отмены ─────────────────────────────────────────────────────
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✖️ Отмена", fmt.Sprintf("%s%d", callbackApartment, aptID)),
	))

	return rows
}
