-- +goose Up
-- Create playlists table
CREATE TABLE IF NOT EXISTS playlists (
    id TEXT NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create playlist_movies junction table
CREATE TABLE IF NOT EXISTS playlist_movies (
    playlist_id TEXT NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    movie_id TEXT NOT NULL,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (playlist_id, movie_id)
);

-- +goose Down
DROP TABLE IF EXISTS playlist_movies;
DROP TABLE IF EXISTS playlists;
