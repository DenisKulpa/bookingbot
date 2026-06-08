CREATE TABLE IF NOT EXISTS users (
    id              SERIAL PRIMARY KEY,
    telegram_id     BIGINT NOT NULL UNIQUE,
    username        TEXT,
    first_name      TEXT,
    last_name       TEXT,
    role            TEXT NOT NULL DEFAULT 'client'
                    CHECK (role IN ('client', 'landlord', 'admin', 'super_admin')),
    is_blocked      INTEGER NOT NULL DEFAULT 0,
    phone           TEXT,
    company_name    TEXT,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Платёжные карты
CREATE TABLE IF NOT EXISTS payment_cards (
    id          SERIAL PRIMARY KEY,
    admin_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label       TEXT NOT NULL,
    card_number TEXT NOT NULL,
    cardholder  TEXT,
    is_active   INTEGER NOT NULL DEFAULT 1,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payment_cards_admin_id ON payment_cards(admin_id);
CREATE INDEX IF NOT EXISTS idx_payment_cards_is_active ON payment_cards(is_active);
