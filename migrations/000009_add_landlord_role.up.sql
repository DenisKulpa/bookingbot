-- Добавляем роль landlord и поля профиля арендодателя

-- Меняем CHECK constraint на role (PostgreSQL требует DROP + ADD)
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('client', 'landlord', 'admin', 'super_admin'));

-- Поля профиля арендодателя (NULL у обычных клиентов)
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone        TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS company_name TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS description  TEXT;
