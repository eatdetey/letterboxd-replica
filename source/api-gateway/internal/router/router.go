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
	authRequired fiber.Handler,
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

	protected := v1.Group("", authRequired)
	protected.Post("/movies/:id/reviews", reviewHandler.AddMovieReview)
	protected.Post("/playlists/search", movieHandler.GetPlaylists)
	protected.Get("/playlists/:id", movieHandler.GetPlaylistByID)
	protected.Post("/playlists", movieHandler.CreatePlaylist)
	protected.Put("/playlists/:id", movieHandler.RenamePlaylist)
	protected.Delete("/playlists/:id", movieHandler.DeletePlaylist)
	protected.Post("/playlists/:id/movies", movieHandler.AddMovieToPlaylist)
	protected.Delete("/playlists/:id/movies/:movie_id", movieHandler.RemoveMovieFromPlaylist)

	app.Use(func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "route not found")
	})
}
