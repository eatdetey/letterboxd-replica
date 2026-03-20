package router

import (
	"github.com/eatdetey/letterboxd-replica/source/api-gateway/internal/handler"
	"github.com/eatdetey/letterboxd-replica/source/api-gateway/internal/middleware"
	"github.com/gofiber/fiber/v3"
)

func Setup(
	app *fiber.App,
	userHandler *handler.UserHandler,
	movieHandler *handler.MovieHandler,
	reviewHandler *handler.ReviewHandler,
) {
	app.Use(middleware.RequestID())

	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})
	app.Get("/swagger", handler.SwaggerUI)
	app.Get("/swagger/", handler.SwaggerUI)
	app.Get("/swagger/openapi.yaml", handler.SwaggerSpec)

	v1 := app.Group("/api/v1")
	v1.Post("/auth/register", userHandler.Register)
	v1.Post("/auth/login", userHandler.Login)
	v1.Post("/auth/refresh", userHandler.Refresh)
	v1.Post("/auth/logout", userHandler.Logout)
	v1.Post("/users/search", userHandler.GetUsers)
	v1.Get("/users/:id", userHandler.GetUserByID)
	v1.Post("/movies/search", movieHandler.GetMovies)
	v1.Get("/movies/:id", movieHandler.GetMovieByID)
	v1.Post("/movies/:id/reviews/search", reviewHandler.GetMovieReviews)
	v1.Get("/movies/:id/reviews/:review_id", reviewHandler.GetMovieReviewByID)
	v1.Post("/movies/:id/reviews", reviewHandler.AddMovieReview)
	v1.Post("/playlists/search", movieHandler.GetPlaylists)
	v1.Get("/playlists/:id", movieHandler.GetPlaylistByID)
	v1.Post("/playlists", movieHandler.CreatePlaylist)
	v1.Put("/playlists/:id", movieHandler.RenamePlaylist)
	v1.Delete("/playlists/:id", movieHandler.DeletePlaylist)
	v1.Post("/playlists/:id/movies", movieHandler.AddMovieToPlaylist)
	v1.Delete("/playlists/:id/movies/:movie_id", movieHandler.RemoveMovieFromPlaylist)

	app.Use(func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "route not found")
	})
}
