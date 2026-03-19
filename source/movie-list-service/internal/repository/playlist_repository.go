package repository

import (
	"context"
	"errors"
	"time"

	pgclient "github.com/eatdetey/letterboxd-replica/source/movie-list-service/internal/transport/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("not found")

type Playlist struct {
	ID          string
	UserID      int64
	Name        string
	MoviesCount int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PlaylistMovie struct {
	PlaylistID string
	MovieID    string
	AddedAt    time.Time
}

type PlaylistRepository interface {
	// Playlist CRUD
	CreatePlaylist(ctx context.Context, playlist Playlist) (*Playlist, error)
	GetPlaylistByID(ctx context.Context, id string) (*Playlist, error)
	GetPlaylistsByUserID(ctx context.Context, userID int64, limit, offset int32) ([]Playlist, error)
	UpdatePlaylistName(ctx context.Context, id, newName string) (*Playlist, error)
	DeletePlaylist(ctx context.Context, id string) error

	// Playlist movie operations
	AddMovieToPlaylist(ctx context.Context, playlistID, movieID string) error
	RemoveMovieFromPlaylist(ctx context.Context, playlistID, movieID string) error
	GetPlaylistMovies(ctx context.Context, playlistID string, limit, offset int32) ([]string, int32, error)
	FilterMoviesByPlaylist(ctx context.Context, playlistID string, candidateMovieIDs []string) ([]string, error)
	GetPlaylistsByMovieID(ctx context.Context, userID int64, movieID string) ([]Playlist, error)
	GetPlaylistMoviesCount(ctx context.Context, playlistID string) (int32, error)
}

type playlistRepository struct {
	pg *pgclient.PostgresClient
}

func NewPlaylistRepository(pg *pgclient.PostgresClient) PlaylistRepository {
	return &playlistRepository{pg: pg}
}

func (r *playlistRepository) CreatePlaylist(ctx context.Context, playlist Playlist) (*Playlist, error) {
	if playlist.UserID == 0 {
		return nil, errors.New("user_id is required")
	}
	if playlist.Name == "" {
		return nil, errors.New("name is required")
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	id := uuid.New().String()
	now := time.Now().UTC()

	row := conn.QueryRow(ctx, createPlaylistQuery, id, playlist.UserID, playlist.Name, now, now)

	var created Playlist
	if err := row.Scan(&created.ID, &created.UserID, &created.Name, &created.CreatedAt, &created.UpdatedAt); err != nil {
		return nil, err
	}
	created.MoviesCount = 0

	return &created, nil
}

func (r *playlistRepository) GetPlaylistByID(ctx context.Context, id string) (*Playlist, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, getPlaylistByIDQuery, id)

	playlist, err := scanPlaylistRow(row)
	if err != nil {
		return nil, err
	}

	// Get movies count
	count, err := r.GetPlaylistMoviesCount(ctx, id)
	if err != nil {
		return nil, err
	}
	playlist.MoviesCount = count

	return playlist, nil
}

func (r *playlistRepository) GetPlaylistsByUserID(ctx context.Context, userID int64, limit, offset int32) ([]Playlist, error) {
	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, getPlaylistsByUserIDQuery, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	playlists := make([]Playlist, 0)
	for rows.Next() {
		playlist, err := scanPlaylistRow(rows)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, *playlist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Get movies count for each playlist
	for i := range playlists {
		count, err := r.GetPlaylistMoviesCount(ctx, playlists[i].ID)
		if err != nil {
			return nil, err
		}
		playlists[i].MoviesCount = count
	}

	return playlists, nil
}

func (r *playlistRepository) UpdatePlaylistName(ctx context.Context, id, newName string) (*Playlist, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if newName == "" {
		return nil, errors.New("new name is required")
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	now := time.Now().UTC()
	row := conn.QueryRow(ctx, updatePlaylistNameQuery, id, newName, now)

	playlist, err := scanPlaylistRow(row)
	if err != nil {
		return nil, err
	}

	count, err := r.GetPlaylistMoviesCount(ctx, id)
	if err != nil {
		return nil, err
	}
	playlist.MoviesCount = count

	return playlist, nil
}

func (r *playlistRepository) DeletePlaylist(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tag, err := conn.Exec(ctx, deletePlaylistQuery, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *playlistRepository) AddMovieToPlaylist(ctx context.Context, playlistID, movieID string) error {
	if playlistID == "" {
		return errors.New("playlist_id is required")
	}
	if movieID == "" {
		return errors.New("movie_id is required")
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	now := time.Now().UTC()
	_, err = conn.Exec(ctx, addMovieToPlaylistQuery, playlistID, movieID, now)
	return err
}

func (r *playlistRepository) RemoveMovieFromPlaylist(ctx context.Context, playlistID, movieID string) error {
	if playlistID == "" {
		return errors.New("playlist_id is required")
	}
	if movieID == "" {
		return errors.New("movie_id is required")
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tag, err := conn.Exec(ctx, removeMovieFromPlaylistQuery, playlistID, movieID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *playlistRepository) GetPlaylistMovies(ctx context.Context, playlistID string, limit, offset int32) ([]string, int32, error) {
	if playlistID == "" {
		return nil, 0, errors.New("playlist_id is required")
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Release()

	// Get total count
	var total int32
	err = conn.QueryRow(ctx, getPlaylistMoviesCountQuery, playlistID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := conn.Query(ctx, getPlaylistMoviesQuery, playlistID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	movieIDs := make([]string, 0)
	for rows.Next() {
		var movieID string
		if err := rows.Scan(&movieID); err != nil {
			return nil, 0, err
		}
		movieIDs = append(movieIDs, movieID)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return movieIDs, total, nil
}

func (r *playlistRepository) FilterMoviesByPlaylist(ctx context.Context, playlistID string, candidateMovieIDs []string) ([]string, error) {
	if playlistID == "" {
		return nil, errors.New("playlist_id is required")
	}
	if len(candidateMovieIDs) == 0 {
		return []string{}, nil
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, filterMoviesByPlaylistQuery, playlistID, candidateMovieIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	filteredIDs := make([]string, 0)
	for rows.Next() {
		var movieID string
		if err := rows.Scan(&movieID); err != nil {
			return nil, err
		}
		filteredIDs = append(filteredIDs, movieID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return filteredIDs, nil
}

func (r *playlistRepository) GetPlaylistsByMovieID(ctx context.Context, userID int64, movieID string) ([]Playlist, error) {
	if movieID == "" {
		return nil, errors.New("movie_id is required")
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, getPlaylistsByMovieIDQuery, userID, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	playlists := make([]Playlist, 0)
	for rows.Next() {
		playlist, err := scanPlaylistRow(rows)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, *playlist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Get movies count for each playlist
	for i := range playlists {
		count, err := r.GetPlaylistMoviesCount(ctx, playlists[i].ID)
		if err != nil {
			return nil, err
		}
		playlists[i].MoviesCount = count
	}

	return playlists, nil
}

func (r *playlistRepository) GetPlaylistMoviesCount(ctx context.Context, playlistID string) (int32, error) {
	if playlistID == "" {
		return 0, errors.New("playlist_id is required")
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	var count int32
	err = conn.QueryRow(ctx, getPlaylistMoviesCountQuery, playlistID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func scanPlaylistRow(row pgx.Row) (*Playlist, error) {
	var playlist Playlist
	if err := row.Scan(&playlist.ID, &playlist.UserID, &playlist.Name, &playlist.CreatedAt, &playlist.UpdatedAt); err != nil {
		return nil, err
	}
	return &playlist, nil
}

const createPlaylistQuery = `
INSERT INTO playlists (id, user_id, name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, name, created_at, updated_at;
`

const getPlaylistByIDQuery = `
SELECT id, user_id, name, created_at, updated_at
FROM playlists
WHERE id = $1;
`

const getPlaylistsByUserIDQuery = `
SELECT id, user_id, name, created_at, updated_at
FROM playlists
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
`

const updatePlaylistNameQuery = `
UPDATE playlists
SET name = $2, updated_at = $3
WHERE id = $1
RETURNING id, user_id, name, created_at, updated_at;
`

const deletePlaylistQuery = `
DELETE FROM playlists WHERE id = $1;
`

const addMovieToPlaylistQuery = `
INSERT INTO playlist_movies (playlist_id, movie_id, added_at)
VALUES ($1, $2, $3)
ON CONFLICT (playlist_id, movie_id) DO NOTHING;
`

const removeMovieFromPlaylistQuery = `
DELETE FROM playlist_movies
WHERE playlist_id = $1 AND movie_id = $2;
`

const getPlaylistMoviesQuery = `
SELECT movie_id
FROM playlist_movies
WHERE playlist_id = $1
ORDER BY added_at DESC
LIMIT $2 OFFSET $3;
`

const getPlaylistMoviesCountQuery = `
SELECT COUNT(*)::int
FROM playlist_movies
WHERE playlist_id = $1;
`

const filterMoviesByPlaylistQuery = `
SELECT movie_id
FROM playlist_movies
WHERE playlist_id = $1 AND movie_id = ANY($2::text[]);
`

const getPlaylistsByMovieIDQuery = `
SELECT p.id, p.user_id, p.name, p.created_at, p.updated_at
FROM playlists p
INNER JOIN playlist_movies pm ON p.id = pm.playlist_id
WHERE p.user_id = $1 AND pm.movie_id = $2;
`
