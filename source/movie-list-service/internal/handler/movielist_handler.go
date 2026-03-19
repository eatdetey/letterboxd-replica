package handler

import (
	"context"
	"strconv"

	"github.com/eatdetey/letterboxd-replica/source/go-common/pkg/ctxmetadata"
	movielistpb "github.com/eatdetey/letterboxd-replica/source/movie-list-service/gen/go/movielist/v1"
	"github.com/eatdetey/letterboxd-replica/source/movie-list-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MovieListHandler struct {
	movielistpb.UnimplementedMovieListServiceServer

	service *service.MovieListService
}

func NewMovieListHandler(service *service.MovieListService) *MovieListHandler {
	return &MovieListHandler{
		service: service,
	}
}

func (h *MovieListHandler) FilterMoviesByPlaylist(ctx context.Context, req *movielistpb.FilterMoviesByPlaylistRequest) (*movielistpb.FilterMoviesByPlaylistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	userID, err := h.userIDFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	filteredIDs, err := h.service.FilterMoviesByPlaylist(ctx, userID, req.PlaylistId, req.CandidateMovieIds)
	if err != nil {
		return nil, err
	}

	return &movielistpb.FilterMoviesByPlaylistResponse{
		MovieIds: filteredIDs,
	}, nil
}

func (h *MovieListHandler) GetPlaylistsForMovie(ctx context.Context, req *movielistpb.GetPlaylistsForMovieRequest) (*movielistpb.GetPlaylistsForMovieResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "movie_id is required")
	}

	userID, err := h.userIDFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	playlists, err := h.service.GetPlaylistsForMovie(ctx, userID, req.MovieId)
	if err != nil {
		return nil, err
	}

	pbPlaylists := make([]*movielistpb.PlaylistInfo, 0, len(playlists))
	for _, pl := range playlists {
		pbPlaylists = append(pbPlaylists, &movielistpb.PlaylistInfo{
			Id:          pl.ID,
			Name:        pl.Name,
			MoviesCount: pl.MoviesCount,
		})
	}

	return &movielistpb.GetPlaylistsForMovieResponse{
		Playlists: pbPlaylists,
	}, nil
}

func (h *MovieListHandler) GetPlaylistsForUser(ctx context.Context, req *movielistpb.GetPlaylistsForUserRequest) (*movielistpb.GetPlaylistsForUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	userID, err := h.userIDFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	playlists, err := h.service.GetPlaylistsForUser(ctx, userID, 0, 0)
	if err != nil {
		return nil, err
	}

	pbPlaylists := make([]*movielistpb.PlaylistInfo, 0, len(playlists))
	for _, pl := range playlists {
		pbPlaylists = append(pbPlaylists, &movielistpb.PlaylistInfo{
			Id:          pl.ID,
			Name:        pl.Name,
			MoviesCount: pl.MoviesCount,
		})
	}

	return &movielistpb.GetPlaylistsForUserResponse{
		Playlists: pbPlaylists,
	}, nil
}

func (h *MovieListHandler) CreatePlaylist(ctx context.Context, req *movielistpb.CreatePlaylistRequest) (*movielistpb.CreatePlaylistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	userID, err := h.userIDFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	playlist, err := h.service.CreatePlaylist(ctx, userID, req.Name)
	if err != nil {
		return nil, err
	}

	return &movielistpb.CreatePlaylistResponse{
		Id:          playlist.ID,
		UserId:      formatUserID(userID),
		Name:        playlist.Name,
		MoviesCount: playlist.MoviesCount,
	}, nil
}

func (h *MovieListHandler) RenamePlaylist(ctx context.Context, req *movielistpb.RenamePlaylistRequest) (*movielistpb.RenamePlaylistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.PlaylistId == "" {
		return nil, status.Error(codes.InvalidArgument, "playlist_id is required")
	}
	if req.NewName == "" {
		return nil, status.Error(codes.InvalidArgument, "new_name is required")
	}

	userID, err := h.userIDFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	playlist, err := h.service.RenamePlaylist(ctx, userID, req.PlaylistId, req.NewName)
	if err != nil {
		return nil, err
	}

	return &movielistpb.RenamePlaylistResponse{
		Id:          playlist.ID,
		UserId:      formatUserID(playlist.UserID),
		Name:        playlist.Name,
		MoviesCount: playlist.MoviesCount,
	}, nil
}

func (h *MovieListHandler) DeletePlaylist(ctx context.Context, req *movielistpb.DeletePlaylistRequest) (*movielistpb.DeletePlaylistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.PlaylistId == "" {
		return nil, status.Error(codes.InvalidArgument, "playlist_id is required")
	}

	userID, err := h.userIDFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	err = h.service.DeletePlaylist(ctx, userID, req.PlaylistId)
	if err != nil {
		return nil, err
	}

	return &movielistpb.DeletePlaylistResponse{}, nil
}

func (h *MovieListHandler) AddMovieToPlaylist(ctx context.Context, req *movielistpb.AddMovieToPlaylistRequest) (*movielistpb.AddMovieToPlaylistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.PlaylistId == "" {
		return nil, status.Error(codes.InvalidArgument, "playlist_id is required")
	}
	if req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "movie_id is required")
	}

	userID, err := h.userIDFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	err = h.service.AddMovieToPlaylist(ctx, userID, req.PlaylistId, req.MovieId)
	if err != nil {
		return nil, err
	}

	return &movielistpb.AddMovieToPlaylistResponse{}, nil
}

func (h *MovieListHandler) RemoveMovieFromPlaylist(ctx context.Context, req *movielistpb.RemoveMovieFromPlaylistRequest) (*movielistpb.RemoveMovieFromPlaylistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.PlaylistId == "" {
		return nil, status.Error(codes.InvalidArgument, "playlist_id is required")
	}
	if req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "movie_id is required")
	}

	userID, err := h.userIDFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	err = h.service.RemoveMovieFromPlaylist(ctx, userID, req.PlaylistId, req.MovieId)
	if err != nil {
		return nil, err
	}

	return &movielistpb.RemoveMovieFromPlaylistResponse{}, nil
}

func (h *MovieListHandler) GetPlaylistMovies(ctx context.Context, req *movielistpb.GetPlaylistMoviesRequest) (*movielistpb.GetPlaylistMoviesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.PlaylistId == "" {
		return nil, status.Error(codes.InvalidArgument, "playlist_id is required")
	}

	userID, err := h.userIDFromClaims(ctx)
	if err != nil {
		return nil, err
	}

	movieIDs, total, err := h.service.GetPlaylistMovies(ctx, userID, req.PlaylistId, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	return &movielistpb.GetPlaylistMoviesResponse{
		MovieIds: movieIDs,
		Total:    total,
	}, nil
}

func (h *MovieListHandler) userIDFromClaims(ctx context.Context) (int64, error) {
	claims, ok := ctxmetadata.GetClaimsFromContext(ctx)
	if !ok || claims == nil {
		return 0, status.Error(codes.Unauthenticated, "authorization required")
	}
	if claims.Id == 0 {
		return 0, status.Error(codes.Unauthenticated, "invalid user claims")
	}
	return claims.Id, nil
}

func formatUserID(userID int64) string {
	return strconv.FormatInt(userID, 10)
}
