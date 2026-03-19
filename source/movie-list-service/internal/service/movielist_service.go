package service

import (
	"context"
	"errors"

	"github.com/eatdetey/letterboxd-replica/source/movie-list-service/internal/repository"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MovieListService struct {
	repo repository.PlaylistRepository
	log  *zap.SugaredLogger
}

func NewMovieListService(repo repository.PlaylistRepository, log *zap.SugaredLogger) *MovieListService {
	return &MovieListService{
		repo: repo,
		log:  log,
	}
}

// FilterMoviesByPlaylist filters candidate movie IDs by membership in a playlist.
func (s *MovieListService) FilterMoviesByPlaylist(ctx context.Context, userID int64, playlistID string, candidateMovieIDs []string) ([]string, error) {
	if len(candidateMovieIDs) == 0 {
		return []string{}, nil
	}

	if _, err := s.getOwnedPlaylist(ctx, playlistID, userID); err != nil {
		return nil, err
	}

	filteredIDs, err := s.repo.FilterMoviesByPlaylist(ctx, playlistID, candidateMovieIDs)
	if err != nil {
		s.log.Errorw("service.filter_movies_by_playlist.filter_failed", "err", err)
		return nil, status.Error(codes.Internal, "failed to filter movies")
	}

	return filteredIDs, nil
}

// GetPlaylistsForMovie returns all playlists containing a specific movie for a user.
func (s *MovieListService) GetPlaylistsForMovie(ctx context.Context, userID int64, movieID string) ([]repository.Playlist, error) {
	if movieID == "" {
		return nil, status.Error(codes.InvalidArgument, "movie_id is required")
	}

	playlists, err := s.repo.GetPlaylistsByMovieID(ctx, userID, movieID)
	if err != nil {
		s.log.Errorw("service.get_playlists_for_movie.failed", "err", err)
		return nil, status.Error(codes.Internal, "failed to get playlists")
	}

	return playlists, nil
}

// GetPlaylistsForUser returns all playlists for a user.
func (s *MovieListService) GetPlaylistsForUser(ctx context.Context, userID int64, limit, offset int32) ([]repository.Playlist, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	playlists, err := s.repo.GetPlaylistsByUserID(ctx, userID, limit, offset)
	if err != nil {
		s.log.Errorw("service.get_playlists_for_user.failed", "err", err)
		return nil, status.Error(codes.Internal, "failed to get playlists")
	}

	return playlists, nil
}

// CreatePlaylist creates a new playlist for a user.
func (s *MovieListService) CreatePlaylist(ctx context.Context, userID int64, name string) (*repository.Playlist, error) {
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	playlist, err := s.repo.CreatePlaylist(ctx, repository.Playlist{
		UserID: userID,
		Name:   name,
	})
	if err != nil {
		s.log.Errorw("service.create_playlist.failed", "err", err)
		return nil, status.Error(codes.Internal, "failed to create playlist")
	}

	return playlist, nil
}

// RenamePlaylist renames an existing playlist.
func (s *MovieListService) RenamePlaylist(ctx context.Context, userID int64, playlistID, newName string) (*repository.Playlist, error) {
	if newName == "" {
		return nil, status.Error(codes.InvalidArgument, "new_name is required")
	}
	if _, err := s.getOwnedPlaylist(ctx, playlistID, userID); err != nil {
		return nil, err
	}

	playlist, err := s.repo.UpdatePlaylistName(ctx, playlistID, newName)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.NotFound, "playlist not found")
		}
		s.log.Errorw("service.rename_playlist.failed", "err", err)
		return nil, status.Error(codes.Internal, "failed to rename playlist")
	}

	return playlist, nil
}

// DeletePlaylist deletes a playlist.
func (s *MovieListService) DeletePlaylist(ctx context.Context, userID int64, playlistID string) error {
	if _, err := s.getOwnedPlaylist(ctx, playlistID, userID); err != nil {
		return err
	}

	err := s.repo.DeletePlaylist(ctx, playlistID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, pgx.ErrNoRows) {
			return status.Error(codes.NotFound, "playlist not found")
		}
		s.log.Errorw("service.delete_playlist.failed", "err", err)
		return status.Error(codes.Internal, "failed to delete playlist")
	}

	return nil
}

// AddMovieToPlaylist adds a movie to a playlist.
func (s *MovieListService) AddMovieToPlaylist(ctx context.Context, userID int64, playlistID, movieID string) error {
	if movieID == "" {
		return status.Error(codes.InvalidArgument, "movie_id is required")
	}
	if _, err := s.getOwnedPlaylist(ctx, playlistID, userID); err != nil {
		return err
	}

	err := s.repo.AddMovieToPlaylist(ctx, playlistID, movieID)
	if err != nil {
		s.log.Errorw("service.add_movie_to_playlist.failed", "err", err)
		return status.Error(codes.Internal, "failed to add movie to playlist")
	}

	return nil
}

// RemoveMovieFromPlaylist removes a movie from a playlist.
func (s *MovieListService) RemoveMovieFromPlaylist(ctx context.Context, userID int64, playlistID, movieID string) error {
	if movieID == "" {
		return status.Error(codes.InvalidArgument, "movie_id is required")
	}
	if _, err := s.getOwnedPlaylist(ctx, playlistID, userID); err != nil {
		return err
	}

	err := s.repo.RemoveMovieFromPlaylist(ctx, playlistID, movieID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return status.Error(codes.NotFound, "playlist or movie not found")
		}
		s.log.Errorw("service.remove_movie_from_playlist.failed", "err", err)
		return status.Error(codes.Internal, "failed to remove movie from playlist")
	}

	return nil
}

// GetPlaylistMovies returns all movies in a playlist.
func (s *MovieListService) GetPlaylistMovies(ctx context.Context, userID int64, playlistID string, limit, offset int32) ([]string, int32, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	if _, err := s.getOwnedPlaylist(ctx, playlistID, userID); err != nil {
		return nil, 0, err
	}

	movieIDs, total, err := s.repo.GetPlaylistMovies(ctx, playlistID, limit, offset)
	if err != nil {
		s.log.Errorw("service.get_playlist_movies.failed", "err", err)
		return nil, 0, status.Error(codes.Internal, "failed to get playlist movies")
	}

	return movieIDs, total, nil
}

func (s *MovieListService) getOwnedPlaylist(ctx context.Context, playlistID string, userID int64) (*repository.Playlist, error) {
	if playlistID == "" {
		return nil, status.Error(codes.InvalidArgument, "playlist_id is required")
	}

	playlist, err := s.repo.GetPlaylistByID(ctx, playlistID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "playlist not found")
		}
		s.log.Errorw("service.playlist.fetch_failed", "err", err, "playlist_id", playlistID)
		return nil, status.Error(codes.Internal, "failed to get playlist")
	}

	if playlist.UserID != userID {
		return nil, status.Error(codes.PermissionDenied, "access denied to playlist")
	}

	return playlist, nil
}
