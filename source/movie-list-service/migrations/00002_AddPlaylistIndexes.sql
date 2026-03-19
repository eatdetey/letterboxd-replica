-- +goose Up
-- Create index for faster lookups by user_id
CREATE INDEX IF NOT EXISTS idx_playlists_user_id ON playlists(user_id);

-- Create index for faster movie_id lookups in playlist_movies
CREATE INDEX IF NOT EXISTS idx_playlist_movies_movie_id ON playlist_movies(movie_id);

-- +goose Down
DROP INDEX IF EXISTS idx_playlist_movies_movie_id;
DROP INDEX IF EXISTS idx_playlists_user_id;
