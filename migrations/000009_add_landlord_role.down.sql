ALTER TABLE users DROP COLUMN IF EXISTS description;
ALTER TABLE users DROP COLUMN IF EXISTS company_name;
ALTER TABLE users DROP COLUMN IF EXISTS phone;

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('client', 'admin', 'super_admin'));
