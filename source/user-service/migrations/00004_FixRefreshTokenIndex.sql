-- +goose Up
WITH duplicates AS (
    SELECT
        id,
        ROW_NUMBER() OVER (PARTITION BY refresh_token ORDER BY created_at DESC, id DESC) AS rn
    FROM user_sessions
    WHERE refresh_token IS NOT NULL
)
DELETE FROM user_sessions us
USING duplicates d
WHERE us.id = d.id AND d.rn > 1;

DROP INDEX IF EXISTS user_sessions_refresh_token_uidx;
CREATE UNIQUE INDEX user_sessions_refresh_token_uidx ON user_sessions (refresh_token);

-- +goose Down
DROP INDEX IF EXISTS user_sessions_refresh_token_uidx;
