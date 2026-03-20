package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	moviepb "github.com/eatdetey/letterboxd-replica/source/api-gateway/gen/go/movie/v1"
	"github.com/eatdetey/letterboxd-replica/source/api-gateway/internal/transport/grpcctx"
	"github.com/gofiber/fiber/v3"
)

const (
	defaultMoviesLimit  int32 = 20
	defaultMoviesOffset int32 = 0
)

type MovieHandler struct {
	client moviepb.MovieServiceClient
}

type movieResponse struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	ReleaseYear int32            `json:"release_year"`
	Genres      []string         `json:"genres"`
	PosterURL   string           `json:"poster_url,omitempty"`
	Playlists   []playlistDetail `json:"playlists,omitempty"`
}

type playlistDetail struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type playlistResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MoviesCount int32  `json:"movies_count"`
}

type createPlaylistRequest struct {
	Name string `json:"name"`
}

type renamePlaylistRequest struct {
	Name string `json:"name"`
}

type addMovieToPlaylistRequest struct {
	MovieID string `json:"movie_id"`
}

func NewMovieHandler(client moviepb.MovieServiceClient) *MovieHandler {
	return &MovieHandler{client: client}
}

func (h *MovieHandler) GetMovies(c fiber.Ctx) error {
	limit, err := parseInt32Query(c.Query("limit"), defaultMoviesLimit, "limit")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	offset, err := parseInt32Query(c.Query("offset"), defaultMoviesOffset, "offset")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if limit <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "limit must be positive")
	}
	if offset < 0 {
		return fiber.NewError(fiber.StatusBadRequest, "offset must be >= 0")
	}

	req := &moviepb.GetMoviesRequest{
		Limit:  limit,
		Offset: offset,
	}

	if search := strings.TrimSpace(c.Query("search")); search != "" {
		req.SearchQuery = &search
	}
	if genre := strings.TrimSpace(c.Query("genre")); genre != "" {
		req.Genre = &genre
	}
	if ids := parseStrings(c.Query("ids")); len(ids) > 0 {
		req.Ids = ids
	}
	if playlistID := strings.TrimSpace(c.Query("playlist_id")); playlistID != "" {
		req.PlaylistId = &playlistID
	}
	if enrichRaw := strings.TrimSpace(c.Query("enrich_playlists")); enrichRaw != "" {
		enrichPlaylists, err := strconv.ParseBool(enrichRaw)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "enrich_playlists must be boolean")
		}
		req.EnrichPlaylists = enrichPlaylists
	}

	ctx, cancel := buildGRPCContext(c)
	defer cancel()

	resp, err := h.client.GetMovies(ctx, req)
	if err != nil {
		return grpcErrorToFiber(err)
	}

	items := make([]movieResponse, 0, len(resp.Items))
	for _, item := range resp.Items {
		playlists := make([]playlistDetail, 0, len(item.Playlists))
		for _, pl := range item.Playlists {
			playlists = append(playlists, playlistDetail{
				ID:   pl.Id,
				Name: pl.Name,
			})
		}

		items = append(items, movieResponse{
			ID:          item.Id,
			Title:       item.Title,
			Description: item.Description,
			ReleaseYear: item.ReleaseYear,
			Genres:      item.Genres,
			PosterURL:   item.Poster,
			Playlists:   playlists,
		})
	}

	return c.JSON(fiber.Map{
		"items": items,
		"total": resp.Total,
	})
}

func (h *MovieHandler) GetPlaylists(c fiber.Ctx) error {
	ctx, cancel := buildGRPCContext(c)
	defer cancel()

	resp, err := h.client.GetPlaylistsForUser(ctx, &moviepb.GetPlaylistsForUserRequest{})
	if err != nil {
		return grpcErrorToFiber(err)
	}

	items := make([]playlistResponse, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, playlistResponse{
			ID:          item.Id,
			Name:        item.Name,
			MoviesCount: item.MoviesCount,
		})
	}

	return c.JSON(fiber.Map{"items": items})
}

func (h *MovieHandler) CreatePlaylist(c fiber.Ctx) error {
	var reqBody createPlaylistRequest
	if err := c.Bind().Body(&reqBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	ctx, cancel := buildGRPCContext(c)
	defer cancel()

	resp, err := h.client.CreatePlaylist(ctx, &moviepb.CreatePlaylistRequest{
		Name: reqBody.Name,
	})
	if err != nil {
		return grpcErrorToFiber(err)
	}

	playlist := resp.GetPlaylist()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"playlist": playlistResponse{
			ID:          playlist.GetId(),
			Name:        playlist.GetName(),
			MoviesCount: playlist.GetMoviesCount(),
		},
	})
}

func (h *MovieHandler) RenamePlaylist(c fiber.Ctx) error {
	playlistID := strings.TrimSpace(c.Params("id"))
	if playlistID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "playlist id is required")
	}

	var reqBody renamePlaylistRequest
	if err := c.Bind().Body(&reqBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	ctx, cancel := buildGRPCContext(c)
	defer cancel()

	resp, err := h.client.RenamePlaylist(ctx, &moviepb.RenamePlaylistRequest{
		PlaylistId: playlistID,
		NewName:    reqBody.Name,
	})
	if err != nil {
		return grpcErrorToFiber(err)
	}

	playlist := resp.GetPlaylist()
	return c.JSON(fiber.Map{
		"playlist": playlistResponse{
			ID:          playlist.GetId(),
			Name:        playlist.GetName(),
			MoviesCount: playlist.GetMoviesCount(),
		},
	})
}

func (h *MovieHandler) DeletePlaylist(c fiber.Ctx) error {
	playlistID := strings.TrimSpace(c.Params("id"))
	if playlistID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "playlist id is required")
	}

	ctx, cancel := buildGRPCContext(c)
	defer cancel()

	_, err := h.client.DeletePlaylist(ctx, &moviepb.DeletePlaylistRequest{
		PlaylistId: playlistID,
	})
	if err != nil {
		return grpcErrorToFiber(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *MovieHandler) AddMovieToPlaylist(c fiber.Ctx) error {
	playlistID := strings.TrimSpace(c.Params("id"))
	if playlistID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "playlist id is required")
	}

	var reqBody addMovieToPlaylistRequest
	if err := c.Bind().Body(&reqBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	ctx, cancel := buildGRPCContext(c)
	defer cancel()

	_, err := h.client.AddMovieToPlaylist(ctx, &moviepb.AddMovieToPlaylistRequest{
		PlaylistId: playlistID,
		MovieId:    reqBody.MovieID,
	})
	if err != nil {
		return grpcErrorToFiber(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *MovieHandler) RemoveMovieFromPlaylist(c fiber.Ctx) error {
	playlistID := strings.TrimSpace(c.Params("id"))
	if playlistID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "playlist id is required")
	}
	movieID := strings.TrimSpace(c.Params("movie_id"))
	if movieID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "movie_id is required")
	}

	ctx, cancel := buildGRPCContext(c)
	defer cancel()

	_, err := h.client.RemoveMovieFromPlaylist(ctx, &moviepb.RemoveMovieFromPlaylistRequest{
		PlaylistId: playlistID,
		MovieId:    movieID,
	})
	if err != nil {
		return grpcErrorToFiber(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func buildGRPCContext(c fiber.Ctx) (context.Context, context.CancelFunc) {
	reqID, _ := c.Locals("request_id").(string)
	ctx := grpcctx.FromFiber(c, reqID)
	return context.WithTimeout(ctx, 5*time.Second)
}

func parseInt32Query(raw string, fallback int32, fieldName string) (int32, error) {
	if strings.TrimSpace(raw) == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be integer", fieldName)
	}
	return int32(value), nil
}
