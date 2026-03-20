package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestRequestID_NoAuthorizationOnPlaylistPath_DoesNotPanic(t *testing.T) {
	app := fiber.New()
	app.Use(RequestID())
	app.Get("/playlists", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest("GET", "/playlists", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("unexpected status code: got %d, want %d", resp.StatusCode, fiber.StatusNoContent)
	}
	if reqID := resp.Header.Get("X-Request-Id"); reqID == "" {
		t.Fatal("expected X-Request-Id header to be set")
	}
}
