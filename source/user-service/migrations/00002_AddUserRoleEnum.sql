-- +goose Up
CREATE TYPE user_role AS ENUM ('user', 'admin');

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role user_role NOT NULL DEFAULT 'user';

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS role;
DROP TYPE IF EXISTS user_role;
