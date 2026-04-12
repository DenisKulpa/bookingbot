CREATE TABLE IF NOT EXISTS users (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id     INTEGER NOT NULL UNIQUE,
    username        TEXT,
    first_name      TEXT,
    last_name       TEXT,
    role            TEXT NOT NULL DEFAULT 'client' CHECK (role IN ('client', 'admin', 'super_admin')),
    is_blocked      INTEGER NOT NULL DEFAULT 0,
    created_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at      DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Seed owner account used by apartment seed data
INSERT OR IGNORE INTO users (telegram_id, username, first_name, last_name, role)
VALUES (100000001, 'booking_admin', 'Администратор', 'Системный', 'admin');
