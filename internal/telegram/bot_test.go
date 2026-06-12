package telegram

import "testing"

func TestToggleFilterAndActiveFilters(t *testing.T) {
	b := &Bot{sessions: make(map[int64]*session)}
	chatID := int64(42)

	b.toggleFilter(chatID, "zone_gagarin_plaza")
	filters := b.activeFilters(chatID)
	if !filters["zone_gagarin_plaza"] {
		t.Fatalf("expected filter to be active after first toggle")
	}

	b.toggleFilter(chatID, "zone_gagarin_plaza")
	filters = b.activeFilters(chatID)
	if filters["zone_gagarin_plaza"] {
		t.Fatalf("expected filter to be inactive after second toggle")
	}
}

func TestResetFilters(t *testing.T) {
	b := &Bot{sessions: make(map[int64]*session)}
	chatID := int64(7)

	b.toggleFilter(chatID, "type_studio")
	b.resetFilters(chatID)

	filters := b.activeFilters(chatID)
	if len(filters) != 0 {
		t.Fatalf("expected empty filters after reset, got %d", len(filters))
	}
}

func TestFindCategoryByFilterCode(t *testing.T) {
	b := &Bot{}
	if got := b.findCategoryByFilterCode("type_studio"); got != "apartment_type" {
		t.Fatalf("expected apartment_type, got %q", got)
	}
	if got := b.findCategoryByFilterCode("unknown_code"); got != "" {
		t.Fatalf("expected empty category for unknown code, got %q", got)
	}
}

func TestFilterListText(t *testing.T) {
	b := &Bot{sessions: make(map[int64]*session)}
	b.sessions[0] = &session{filters: make(map[string]bool), city: "Одесса"}
	if got := b.filterListText(0); got != "🏖 *Поиск жилья — Одесса*\n\nВыберите категорию:" {
		t.Fatalf("unexpected text without filters: %q", got)
	}
	b.sessions[0].filters["a"] = true
	if got := b.filterListText(0); got != "🏖 *Поиск жилья — Одесса*\n\nВыбрано фильтров: *1*\nВыберите категорию:" {
		t.Fatalf("unexpected text with filters: %q", got)
	}
}
