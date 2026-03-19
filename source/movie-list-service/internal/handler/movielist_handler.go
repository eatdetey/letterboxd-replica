package handler

import (
	"context"
	"strconv"

	"github.com/eatdetey/letterboxd-replica/source/movie-list-service/internal/service"
	movielistpb "github.com/eatdetey/letterboxd-replica/source/movie-list-service/gen/go/movielist/v1"
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

func (h *MovieListHandler) FilterMoviesByPlaylist(ctx context.Context, req *movielistpb.FilterRequest) (*movielistpb.MovieIdsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	filteredIDs, err := h.service.FilterMoviesByPlaylist(ctx, userID, req.PlaylistId, req.CandidateMovieIds)
	if err != nil {
		return nil, err
	}

	return &movielistpb.MovieIdsResponse{
		MovieIds: filteredIDs,
	}, nil
}

func (h *MovieListHandler) GetPlaylistsForMovie(ctx context.Context, req *movielistpb.MoviePlaylistRequest) (*movielistpb.PlaylistsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	playlists, err := h.service.GetPlaylistsForMovie(ctx, userID, req.MovieId)
	if err != nil {
		return nil, err
	}

	pbPlaylists := make([]*movielistpb.PlaylistInfo, 0, len(playlists))
	for _, pl := range playlists {
		pbPlaylists = append(pbPlaylists, &movielistpb.PlaylistInfo{
			Id:           pl.ID,
			Name:         pl.Name,
			MoviesCount:  pl.MoviesCount,
		})
	}

	return &movielistpb.PlaylistsResponse{
		Playlists: pbPlaylists,
	}, nil
}

func (h *MovieListHandler) GetPlaylistsForUser(ctx context.Context, req *movielistpb.UserPlaylistsRequest) (*movielistpb.UserPlaylistsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	playlists, err := h.service.GetPlaylistsForUser(ctx, userID, 0, 0)
	if err != nil {
		return nil, err
	}

	pbPlaylists := make([]*movielistpb.PlaylistInfo, 0, len(playlists))
	for _, pl := range playlists {
		pbPlaylists = append(pbPlaylists, &movielistpb.PlaylistInfo{
			Id:           pl.ID,
			Name:         pl.Name,
			MoviesCount:  pl.MoviesCount,
		})
	}

	return &movielistpb.UserPlaylistsResponse{
		Playlists: pbPlaylists,
	}, nil
}

func (h *MovieListHandler) CreatePlaylist(ctx context.Context, req *movielistpb.CreatePlaylistRequest) (*movielistpb.PlaylistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	playlist, err := h.service.CreatePlaylist(ctx, userID, req.Name)
	if err != nil {
		return nil, err
	}

	return &movielistpb.PlaylistResponse{
		Id:          playlist.ID,
		UserId:      req.UserId,
		Name:        playlist.Name,
		MoviesCount: playlist.MoviesCount,
	}, nil
}

func (h *MovieListHandler) RenamePlaylist(ctx context.Context, req *movielistpb.RenamePlaylistRequest) (*movielistpb.PlaylistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.PlaylistId == "" {
		return nil, status.Error(codes.InvalidArgument, "playlist_id is required")
	}

	playlist, err := h.service.RenamePlaylist(ctx, req.PlaylistId, req.NewName)
	if err != nil {
		return nil, err
	}

	return &movielistpb.PlaylistResponse{
		Id:          playlist.ID,
		UserId:      strconv.FormatInt(playlist.UserID, 10),
		Name:        playlist.Name,
		MoviesCount: playlist.MoviesCount,
	}, nil
}

func (h *MovieListHandler) DeletePlaylist(ctx context.Context, req *movielistpb.DeletePlaylistRequest) (*movielistpb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.PlaylistId == "" {
		return nil, status.Error(codes.InvalidArgument, "playlist_id is required")
	}

	err := h.service.DeletePlaylist(ctx, req.PlaylistId)
	if err != nil {
		return nil, err
	}

	return &movielistpb.Empty{}, nil
}

func (h *MovieListHandler) AddMovieToPlaylist(ctx context.Context, req *movielistpb.AddMovieRequest) (*movielistpb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.PlaylistId == "" {
		return nil, status.Error(codes.InvalidArgument, "playlist_id is required")
	}
	if req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "movie_id is required")
	}

	err := h.service.AddMovieToPlaylist(ctx, req.PlaylistId, req.MovieId)
	if err != nil {
		return nil, err
	}

	return &movielistpb.Empty{}, nil
}

func (h *MovieListHandler) RemoveMovieFromPlaylist(ctx context.Context, req *movielistpb.RemoveMovieRequest) (*movielistpb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.PlaylistId == "" {
		return nil, status.Error(codes.InvalidArgument, "playlist_id is required")
	}
	if req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "movie_id is required")
	}

	err := h.service.RemoveMovieFromPlaylist(ctx, req.PlaylistId, req.MovieId)
	if err != nil {
		return nil, err
	}

	return &movielistpb.Empty{}, nil
}

func (h *MovieListHandler) GetPlaylistMovies(ctx context.Context, req *movielistpb.GetPlaylistMoviesRequest) (*movielistpb.PlaylistMoviesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.PlaylistId == "" {
		return nil, status.Error(codes.InvalidArgument, "playlist_id is required")
	}

	movieIDs, total, err := h.service.GetPlaylistMovies(ctx, req.PlaylistId, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	return &movielistpb.PlaylistMoviesResponse{
		MovieIds: movieIDs,
		Total:    total,
	}, nil
}
