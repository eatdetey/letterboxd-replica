package router

import (
	"github.com/eatdetey/letterboxd-replica/source/api-gateway/internal/handler"
	"github.com/eatdetey/letterboxd-replica/source/api-gateway/internal/middleware"
	"github.com/gofiber/fiber/v3"
)

func Setup(app *fiber.App, userHandler *handler.UserHandler, movieHandler *handler.MovieHandler) {
	app.Use(middleware.RequestID())

	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	v1 := app.Group("/api/v1")
	v1.Post("/auth/register", userHandler.Register)
	v1.Post("/auth/login", userHandler.Login)
	v1.Post("/auth/refresh", userHandler.Refresh)
	v1.Get("/users", userHandler.GetUsers)
	v1.Get("/movies", movieHandler.GetMovies)

	app.Use(func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "route not found")
	})
}
