CREATE TABLE IF NOT EXISTS payment_cards (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    admin_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label       TEXT NOT NULL,           -- например "Приватбанк Visa"
    card_number TEXT NOT NULL,           -- последние 4 цифры или полный номер
    cardholder  TEXT,
    is_active   INTEGER NOT NULL DEFAULT 1,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_payment_cards_admin_id ON payment_cards(admin_id);
CREATE INDEX IF NOT EXISTS idx_payment_cards_is_active ON payment_cards(is_active);
