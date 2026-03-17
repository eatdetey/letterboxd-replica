-- +goose Up
CREATE UNIQUE INDEX IF NOT EXISTS user_sessions_refresh_token_uidx ON user_sessions (refresh_token);

-- +goose Down
DROP INDEX IF EXISTS user_sessions_refresh_token_uidx;
